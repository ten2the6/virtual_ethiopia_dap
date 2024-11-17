package blockchain

import (
    "crypto/sha256"
    "encoding/hex"    
    "errors"
    "sync"
    "time"
)

type ElectionStatus int

const (
    NotStarted ElectionStatus = iota
    InProgress
    Completed
    Cancelled
)

// Election represents a presidential election
type Election struct {
    ID            string         `json:"id"`
    Name          string         `json:"name"`
    StartDate     int64          `json:"startDate"`
    EndDate       int64          `json:"endDate"`
    Status        ElectionStatus `json:"status"`
    Candidates    []Candidate    `json:"candidates"`
    Votes         map[string]string  `json:"votes"`  // CitizenID -> CandidateID
    Winner        *Candidate     `json:"winner,omitempty"`
}

// Candidate represents a presidential candidate
type Candidate struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    PublicKey   string    `json:"publicKey"`
    Platform    string    `json:"platform"`
    VoteCount   int       `json:"voteCount"`
}

// ElectionSystem manages the election process
type ElectionSystem struct {
    currentElection *Election
    pastElections   []*Election
    citizenRegistry *CitizenRegistry
    mu             sync.RWMutex
}

// NewElectionSystem creates a new election system
func NewElectionSystem(registry *CitizenRegistry) *ElectionSystem {
    return &ElectionSystem{
        pastElections:   make([]*Election, 0),
        citizenRegistry: registry,
    }
}

// StartElection initiates a new presidential election
func (es *ElectionSystem) StartElection(name string, durationDays int) error {
    es.mu.Lock()
    defer es.mu.Unlock()

    if es.currentElection != nil && es.currentElection.Status == InProgress {
        return errors.New("an election is already in progress")
    }

    startDate := time.Now().Unix()
    endDate := time.Now().AddDate(0, 0, durationDays).Unix()

    es.currentElection = &Election{
        ID:         generateElectionID(),
        Name:       name,
        StartDate:  startDate,
        EndDate:    endDate,
        Status:     InProgress,
        Candidates: make([]Candidate, 0),
        Votes:      make(map[string]string),
    }

    return nil
}

// GetCurrentElection returns the current election
func (es *ElectionSystem) GetCurrentElection() *Election {
    es.mu.RLock()
    defer es.mu.RUnlock()
    return es.currentElection
}

// RegisterCandidate registers a new presidential candidate
func (es *ElectionSystem) RegisterCandidate(name, publicKey, platform string) error {
    es.mu.Lock()
    defer es.mu.Unlock()

    if es.currentElection == nil || es.currentElection.Status != InProgress {
        return errors.New("no active election")
    }

    if !es.citizenRegistry.IsCitizen(publicKey) {
        return errors.New("candidate must be an approved citizen")
    }

    candidate := Candidate{
        ID:        generateCandidateID(name, publicKey),
        Name:      name,
        PublicKey: publicKey,
        Platform:  platform,
    }

    es.currentElection.Candidates = append(es.currentElection.Candidates, candidate)
    return nil
}

// CastVote records a citizen's vote
func (es *ElectionSystem) CastVote(citizenPublicKey, candidateID string) error {
    es.mu.Lock()
    defer es.mu.Unlock()

    if es.currentElection == nil || es.currentElection.Status != InProgress {
        return errors.New("no active election")
    }

    if !es.citizenRegistry.IsCitizen(citizenPublicKey) {
        return errors.New("voter must be an approved citizen")
    }

    if _, voted := es.currentElection.Votes[citizenPublicKey]; voted {
        return errors.New("citizen has already voted")
    }

    candidateExists := false
    for _, candidate := range es.currentElection.Candidates {
        if candidate.ID == candidateID {
            candidateExists = true
            break
        }
    }
    if !candidateExists {
        return errors.New("invalid candidate")
    }

    es.currentElection.Votes[citizenPublicKey] = candidateID
    return nil
}

// EndElection concludes the current election and determines the winner
func (es *ElectionSystem) EndElection() error {
    es.mu.Lock()
    defer es.mu.Unlock()

    if es.currentElection == nil || es.currentElection.Status != InProgress {
        return errors.New("no active election")
    }

    // Count votes
    voteCounts := make(map[string]int)
    for _, candidateID := range es.currentElection.Votes {
        voteCounts[candidateID]++
    }

    // Find winner
    var winningCandidate *Candidate
    maxVotes := 0
    for i, candidate := range es.currentElection.Candidates {
        votes := voteCounts[candidate.ID]
        es.currentElection.Candidates[i].VoteCount = votes
        if votes > maxVotes {
            maxVotes = votes
            winningCandidate = &es.currentElection.Candidates[i]
        }
    }

    es.currentElection.Status = Completed
    es.currentElection.Winner = winningCandidate
    es.pastElections = append(es.pastElections, es.currentElection)
    es.currentElection = nil

    return nil
}

func generateElectionID() string {
    h := sha256.New()
    h.Write([]byte(time.Now().String()))
    return hex.EncodeToString(h.Sum(nil))
}

func generateCandidateID(name, publicKey string) string {
    h := sha256.New()
    h.Write([]byte(name + publicKey))
    return hex.EncodeToString(h.Sum(nil))
}
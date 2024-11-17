package blockchain

import (
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "sync"
    "time"
)

type CitizenStatus int

const (
    Pending CitizenStatus = iota
    Approved
    Rejected
)

// Citizen represents a citizen of the virtual nation
type Citizen struct {
    ID            string         `json:"id"`
    PublicKey     string         `json:"publicKey"`
    Name          string         `json:"name"`
    DateOfBirth   string         `json:"dateOfBirth"`
    RegisterDate  int64          `json:"registerDate"`
    Status        CitizenStatus  `json:"status"`
    ApprovedBy    string         `json:"approvedBy,omitempty"`
    ApprovalDate  int64          `json:"approvalDate,omitempty"`
}

// CitizenRegistry manages citizen registration
type CitizenRegistry struct {
    citizens map[string]*Citizen  // PublicKey -> Citizen
    admins   map[string]bool     // PublicKey -> isAdmin
    mu       sync.RWMutex
}

// NewCitizenRegistry creates a new citizen registry
func NewCitizenRegistry() *CitizenRegistry {
    registry := &CitizenRegistry{
        citizens: make(map[string]*Citizen),
        admins:   make(map[string]bool),
    }
    // Add initial admin
    registry.admins["GENESIS_ADMIN"] = true
    return registry
}

// RegisterCitizen creates a new citizen registration request
func (cr *CitizenRegistry) RegisterCitizen(name, dateOfBirth, publicKey string) (*Citizen, error) {
    cr.mu.Lock()
    defer cr.mu.Unlock()

    if _, exists := cr.citizens[publicKey]; exists {
        return nil, errors.New("citizen already registered")
    }

    id := generateCitizenID(name, publicKey)
    citizen := &Citizen{
        ID:           id,
        PublicKey:    publicKey,
        Name:         name,
        DateOfBirth:  dateOfBirth,
        RegisterDate: time.Now().Unix(),
        Status:       Pending,
    }

    cr.citizens[publicKey] = citizen
    return citizen, nil
}

// ApproveCitizen approves a citizen registration
func (cr *CitizenRegistry) ApproveCitizen(citizenID, approverKey string) error {
    cr.mu.Lock()
    defer cr.mu.Unlock()

    if !cr.admins[approverKey] {
        return errors.New("not authorized to approve citizens")
    }

    var targetCitizen *Citizen
    for _, citizen := range cr.citizens {
        if citizen.ID == citizenID {
            targetCitizen = citizen
            break
        }
    }

    if targetCitizen == nil {
        return errors.New("citizen not found")
    }

    if targetCitizen.Status != Pending {
        return errors.New("citizen already processed")
    }

    targetCitizen.Status = Approved
    targetCitizen.ApprovedBy = approverKey
    targetCitizen.ApprovalDate = time.Now().Unix()
    return nil
}

// GetCitizen returns a citizen by their public key
func (cr *CitizenRegistry) GetCitizen(publicKey string) (*Citizen, bool) {
    cr.mu.RLock()
    defer cr.mu.RUnlock()
    
    citizen, exists := cr.citizens[publicKey]
    return citizen, exists
}

// GetAllCitizens returns all registered citizens
func (cr *CitizenRegistry) GetAllCitizens() []*Citizen {
    cr.mu.RLock()
    defer cr.mu.RUnlock()
    
    citizens := make([]*Citizen, 0, len(cr.citizens))
    for _, citizen := range cr.citizens {
        citizens = append(citizens, citizen)
    }
    return citizens
}

// IsCitizen checks if a public key belongs to an approved citizen
func (cr *CitizenRegistry) IsCitizen(publicKey string) bool {
    cr.mu.RLock()
    defer cr.mu.RUnlock()

    citizen, exists := cr.citizens[publicKey]
    return exists && citizen.Status == Approved
}

func generateCitizenID(name, publicKey string) string {
    h := sha256.New()
    h.Write([]byte(name + publicKey))
    return hex.EncodeToString(h.Sum(nil))
}
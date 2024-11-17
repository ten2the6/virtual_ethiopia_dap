package blockchain

import (
    "fmt"
    "sync"
)

// Chain represents the blockchain
type Chain struct {
    blocks          []*Block
    mu              sync.RWMutex
    txPool          *TransactionPool
    citizenRegistry *CitizenRegistry
    electionSystem  *ElectionSystem
}

// NewChain creates a new blockchain
func NewChain() *Chain {
    registry := NewCitizenRegistry()
    chain := &Chain{
        blocks:          make([]*Block, 0),
        txPool:          NewTransactionPool(),
        citizenRegistry: registry,
        electionSystem:  NewElectionSystem(registry),
    }
    chain.addGenesisBlock()
    return chain
}

// addGenesisBlock creates and adds the genesis block
func (c *Chain) addGenesisBlock() {
    genesisBlock := NewBlock(0, []Transaction{}, []byte("0"))
    c.blocks = append(c.blocks, genesisBlock)
}

// AddBlock creates and adds a new block
func (c *Chain) AddBlock() error {
    c.mu.Lock()
    defer c.mu.Unlock()

    if len(c.blocks) == 0 {
        return fmt.Errorf("blockchain not initialized")
    }

    prevBlock := c.blocks[len(c.blocks)-1]
    transactions := make([]Transaction, 0)
    for _, tx := range c.txPool.GetAllTransactions() {
        transactions = append(transactions, *tx)
    }

    newBlock := NewBlock(
        prevBlock.Index+1,
        transactions,
        prevBlock.Hash,
    )

    c.blocks = append(c.blocks, newBlock)
    c.txPool.Clear() // Clear the transaction pool
    return nil
}

// AddTransaction adds a transaction to the pool
func (c *Chain) AddTransaction(from, to string, amount float64) (*Transaction, error) {
    tx := NewTransaction(from, to, amount)
    if !c.txPool.AddTransaction(tx) {
        return nil, fmt.Errorf("failed to add transaction to pool")
    }
    return tx, nil
}

// GetLatestBlock returns the most recent block
func (c *Chain) GetLatestBlock() (*Block, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    if len(c.blocks) == 0 {
        return nil, fmt.Errorf("chain is empty")
    }
    return c.blocks[len(c.blocks)-1], nil
}

// GetBlocks returns all blocks
func (c *Chain) GetBlocks() []*Block {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.blocks
}

// ValidateChain verifies the integrity of the blockchain
func (c *Chain) ValidateChain() bool {
    c.mu.RLock()
    defer c.mu.RUnlock()

    for i := 1; i < len(c.blocks); i++ {
        currentBlock := c.blocks[i]
        prevBlock := c.blocks[i-1]

        if string(currentBlock.PrevHash) != string(prevBlock.Hash) {
            return false
        }

        if string(currentBlock.Hash) != string(currentBlock.calculateHash()) {
            return false
        }
    }
    return true
}

// GetTransactionPool returns the current transaction pool
func (c *Chain) GetTransactionPool() *TransactionPool {
    return c.txPool
}

// AddCitizenRegistration adds a new citizen registration transaction
func (c *Chain) AddCitizenRegistration(name, dateOfBirth, publicKey string) (*Transaction, error) {
    citizen, err := c.citizenRegistry.RegisterCitizen(name, dateOfBirth, publicKey)
    if err != nil {
        return nil, err
    }

    tx := NewTransaction("SYSTEM", publicKey, 0)
    tx.Data = map[string]interface{}{
        "type":     "CITIZEN_REGISTRATION",
        "citizen":  citizen,
    }
    
    if !c.txPool.AddTransaction(tx) {
        return nil, fmt.Errorf("failed to add citizen registration transaction")
    }
    return tx, nil
}

// ApproveCitizen approves a citizen registration
func (c *Chain) ApproveCitizen(citizenID, approverKey string) (*Transaction, error) {
    if err := c.citizenRegistry.ApproveCitizen(citizenID, approverKey); err != nil {
        return nil, err
    }

    tx := NewTransaction("SYSTEM", citizenID, 0)
    tx.Data = map[string]interface{}{
        "type":        "CITIZEN_APPROVAL",
        "citizenID":   citizenID,
        "approverKey": approverKey,
    }
    
    if !c.txPool.AddTransaction(tx) {
        return nil, fmt.Errorf("failed to add citizen approval transaction")
    }
    return tx, nil
}

// StartElection starts a new presidential election
func (c *Chain) StartElection(name string, durationDays int) (*Transaction, error) {
    if err := c.electionSystem.StartElection(name, durationDays); err != nil {
        return nil, err
    }

    tx := NewTransaction("SYSTEM", "ELECTION", 0)
    tx.Data = map[string]interface{}{
        "type":         "ELECTION_START",
        "name":         name,
        "durationDays": durationDays,
    }
    
    if !c.txPool.AddTransaction(tx) {
        return nil, fmt.Errorf("failed to add election start transaction")
    }
    return tx, nil
}

// RegisterCandidate registers a new presidential candidate
func (c *Chain) RegisterCandidate(name, publicKey, platform string) (*Transaction, error) {
    if err := c.electionSystem.RegisterCandidate(name, publicKey, platform); err != nil {
        return nil, err
    }

    tx := NewTransaction("SYSTEM", publicKey, 0)
    tx.Data = map[string]interface{}{
        "type":     "CANDIDATE_REGISTRATION",
        "name":     name,
        "platform": platform,
    }
    
    if !c.txPool.AddTransaction(tx) {
        return nil, fmt.Errorf("failed to add candidate registration transaction")
    }
    return tx, nil
}

// CastVote records a citizen's vote
func (c *Chain) CastVote(citizenPublicKey, candidateID string) (*Transaction, error) {
    if err := c.electionSystem.CastVote(citizenPublicKey, candidateID); err != nil {
        return nil, err
    }

    tx := NewTransaction(citizenPublicKey, "ELECTION", 0)
    tx.Data = map[string]interface{}{
        "type":        "VOTE_CAST",
        "candidateID": candidateID,
    }
    
    if !c.txPool.AddTransaction(tx) {
        return nil, fmt.Errorf("failed to add vote transaction")
    }
    return tx, nil
}

// GetCitizen returns a citizen by their public key
func (c *Chain) GetCitizen(publicKey string) (*Citizen, bool) {
    return c.citizenRegistry.GetCitizen(publicKey)
}

// GetAllCitizens returns all registered citizens
func (c *Chain) GetAllCitizens() []*Citizen {
    return c.citizenRegistry.GetAllCitizens()
}

// GetCurrentElection returns the current active election
func (c *Chain) GetCurrentElection() *Election {
    return c.electionSystem.GetCurrentElection()
}

// GetCurrentElectionCandidates returns all candidates in the current election
func (c *Chain) GetCurrentElectionCandidates() []Candidate {
    if election := c.electionSystem.GetCurrentElection(); election != nil {
        return election.Candidates
    }
    return []Candidate{}
}

// EndElection ends the current election and determines the winner
func (c *Chain) EndElection() (*Transaction, error) {
    if err := c.electionSystem.EndElection(); err != nil {
        return nil, err
    }

    tx := NewTransaction("SYSTEM", "ELECTION", 0)
    tx.Data = map[string]interface{}{
        "type": "ELECTION_END",
    }
    
    if !c.txPool.AddTransaction(tx) {
        return nil, fmt.Errorf("failed to add election end transaction")
    }
    return tx, nil
}
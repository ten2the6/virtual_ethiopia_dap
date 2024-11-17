package blockchain

import (
    "crypto/sha256"
    "encoding/hex"
    "time"
)

// Transaction represents a blockchain transaction
type Transaction struct {
    ID        string                 `json:"id"`
    From      string                 `json:"from"`
    To        string                 `json:"to"`
    Amount    float64               `json:"amount"`
    Timestamp int64                 `json:"timestamp"`
    Signature string                `json:"signature,omitempty"`
    Data      map[string]interface{} `json:"data,omitempty"`
}

// NewTransaction creates a new transaction
func NewTransaction(from, to string, amount float64) *Transaction {
    tx := &Transaction{
        From:      from,
        To:        to,
        Amount:    amount,
        Timestamp: time.Now().Unix(),
        Data:      make(map[string]interface{}),
    }
    tx.ID = calculateTransactionHash(tx)
    return tx
}

// calculateTransactionHash generates a hash for the transaction
func calculateTransactionHash(tx *Transaction) string {
    record := tx.From + tx.To + string(tx.Timestamp) + string(int64(tx.Amount*100000000))
    h := sha256.New()
    h.Write([]byte(record))
    return hex.EncodeToString(h.Sum(nil))
}

// SignTransaction adds a signature to the transaction (simplified version)
func SignTransaction(tx *Transaction, privateKey string) {
    // In a real implementation, this would use proper cryptographic signing
    h := sha256.New()
    h.Write([]byte(tx.ID + privateKey))
    tx.Signature = hex.EncodeToString(h.Sum(nil))
}

// VerifyTransaction checks if the transaction signature is valid (simplified version)
func VerifyTransaction(tx *Transaction, publicKey string) bool {
    // In a real implementation, this would verify the cryptographic signature
    return tx.Signature != ""
}

// TransactionPool manages pending transactions
type TransactionPool struct {
    transactions map[string]*Transaction
}

// NewTransactionPool creates a new transaction pool
func NewTransactionPool() *TransactionPool {
    return &TransactionPool{
        transactions: make(map[string]*Transaction),
    }
}

// AddTransaction adds a transaction to the pool
func (tp *TransactionPool) AddTransaction(tx *Transaction) bool {
    if tx == nil || tx.ID == "" {
        return false
    }
    
    // Check if transaction already exists
    if _, exists := tp.transactions[tx.ID]; exists {
        return false
    }
    
    tp.transactions[tx.ID] = tx
    return true
}

// GetTransaction retrieves a transaction from the pool by ID
func (tp *TransactionPool) GetTransaction(id string) (*Transaction, bool) {
    tx, exists := tp.transactions[id]
    return tx, exists
}

// GetAllTransactions returns all transactions in the pool
func (tp *TransactionPool) GetAllTransactions() []*Transaction {
    txs := make([]*Transaction, 0, len(tp.transactions))
    for _, tx := range tp.transactions {
        txs = append(txs, tx)
    }
    return txs
}

// RemoveTransaction removes a transaction from the pool
func (tp *TransactionPool) RemoveTransaction(id string) {
    delete(tp.transactions, id)
}

// Clear removes all transactions from the pool
func (tp *TransactionPool) Clear() {
    tp.transactions = make(map[string]*Transaction)
}

// Size returns the number of transactions in the pool
func (tp *TransactionPool) Size() int {
    return len(tp.transactions)
}
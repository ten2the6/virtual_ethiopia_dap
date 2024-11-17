package blockchain

import (
	"fmt"
	"time"
	"crypto/sha256"
	// Removed unused import: encoding/hex
)

type Block struct {
	Index        int64         `json:"index"`
	Timestamp    int64         `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	Hash         []byte        `json:"hash"`
	PrevHash     []byte        `json:"prevHash"`
}

func NewBlock(index int64, transactions []Transaction, prevHash []byte) *Block {
	block := &Block{
		Index:        index,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
		PrevHash:     prevHash,
	}
	block.Hash = block.calculateHash()
	return block
}

func (b *Block) calculateHash() []byte {
	data := fmt.Sprintf("%d%d%v%s", b.Index, b.Timestamp, b.Transactions, b.PrevHash)
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}


package models

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"
)

type Block struct {
	Hash         string                 `json:"hash"`
	PrevHash     string                 `json:"prevHash"`
	Transactions map[string]Transaction `json:"transactions"`
	Nonce        int64                  `json:"nonce"`
	Timestamp    int64                  `json:"timestamp"`
}

func NewBlock(prevHash string, txs map[string]Transaction) *Block {
	b := Block{
		PrevHash:     prevHash,
		Transactions: make(map[string]Transaction),
		Timestamp:    time.Now().Unix(),
	}

	for k, t := range txs {
		b.Transactions[k] = Transaction{
			Id:        t.Id,
			Candidate: t.Candidate,
			Voter:     t.Voter,
			Timestamp: t.Timestamp,
		}
	}

	return &b
}

func (b *Block) CalculateHash(difficulty int) string {
	h := sha256.New()

	secureHash := strings.Repeat("0", difficulty)
	h.Reset()
	h.Write([]byte(b.PrevHash))
	h.Write([]byte(fmt.Sprintf("%+v", b.Transactions)))
	h.Write([]byte(fmt.Sprintf("%d", b.Timestamp)))
	h.Write([]byte(fmt.Sprintf("%d", b.Nonce)))

	hash := fmt.Sprintf("%x", h.Sum(nil))
	// fmt.Printf("%+v\n", b)
	// fmt.Println(hash[0:difficulty], secureHash)
	if hash[0:difficulty] == secureHash {
		return hash
	}

	b.Nonce = b.Nonce + 1
	return b.CalculateHash(difficulty)
}

func (b *Block) Mine(difficulty int) *Block {
	b.Hash = b.CalculateHash(difficulty)
	return b
}

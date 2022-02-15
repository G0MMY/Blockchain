package components

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"
)

type BlockType struct {
	Id        int   `json:"Id" gorm:"primaryKey"`
	Nonce     int   `json:"nonce"`
	Timestamp int64 `json:"timestamp"`
	//Transactions          []*TransactionType
	PreviousHash          []byte `json:"previousHash"`
	CurrentHash           []byte `json:"currentHash"`
	MaxNumberTransactions int    `json:"maxNumberTransactions"`
}

type Block interface {
	CheckBlock() bool
}

func (block *BlockType) CheckBlock() bool {
	if fmt.Sprintf("%x", hash(block.Nonce, block)) == fmt.Sprintf("%x", block.CurrentHash) {
		return true
	}
	return false
}

func CreateBlock(id int, previousHash []byte) *BlockType {
	block := &BlockType{id, 0, time.Now().Unix(), previousHash, []byte{}, 10}
	ProofOfWork(block)
	return block
}

func ProofOfWork(block *BlockType) {
	i := 0
	stringhash := fmt.Sprintf("%x", hash(i, block))
	for stringhash[0:4] != "0000" {
		i += 1
		stringhash = fmt.Sprintf("%x", hash(i, block))
	}
	block.Nonce = i
	block.CurrentHash = hash(i, block)
}

func hash(nonce int, block *BlockType) []byte {
	info := bytes.Join([][]byte{
		[]byte(fmt.Sprintf("%x", nonce)),
		[]byte(fmt.Sprintf("%x", block.Timestamp)),
		block.PreviousHash,
	}, []byte{})
	hash := sha256.Sum256(info)

	return hash[:]
}

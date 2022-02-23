package Controllers

import (
	"blockchain/Models"
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"
)

func CreateBlock(previousHash []byte) *Models.Block {
	block := &Models.Block{Nonce: 0, Timestamp: time.Now().Unix(), Transactions: []Models.Transaction{}, PreviousHash: previousHash, CurrentHash: []byte{}, MaxNumberTransactions: 10}
	ProofOfWork(block)
	return block
}

func ProofOfWork(block *Models.Block) {
	i := 0
	stringhash := fmt.Sprintf("%x", hash(i, block))
	for stringhash[0:4] != "0000" {
		i += 1
		stringhash = fmt.Sprintf("%x", hash(i, block))
	}
	block.Nonce = i
	block.CurrentHash = hash(i, block)
}

func hash(nonce int, block *Models.Block) []byte {
	info := bytes.Join([][]byte{
		[]byte(fmt.Sprintf("%x", nonce)),
		[]byte(fmt.Sprintf("%x", block.Timestamp)),
		block.PreviousHash,
	}, []byte{})
	hash := sha256.Sum256(info)

	return hash[:]
}

func CheckBlock(block *Models.Block) bool {
	if fmt.Sprintf("%s", hash(block.Nonce, block)) == fmt.Sprintf("%s", block.CurrentHash) {
		return true
	}
	return false
}

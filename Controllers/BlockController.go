package Controllers

import (
	"blockchain/Models"
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"
)

func CreateBlock(previousHash []byte, transactions []Models.Transaction) *Models.Block {
	block := &Models.Block{Nonce: 0, Timestamp: time.Now().Unix(), MerkleRoot: []byte{}, PreviousHash: previousHash, Difficulty: 4, Transactions: transactions}
	ProofOfWork(block, transactions)
	return block
}

//****************Implement difficulty***********************
func ProofOfWork(block *Models.Block, transactions []Models.Transaction) {
	i := 0
	stringhash := fmt.Sprintf("%x", hash(i, block, transactions))
	for stringhash[0:4] != "0000" {
		i += 1
		stringhash = fmt.Sprintf("%x", hash(i, block, transactions))
	}
	block.Nonce = i
}

func hash(nonce int, block *Models.Block, transactions []Models.Transaction) []byte {
	info := bytes.Join([][]byte{
		[]byte(fmt.Sprintf("%x", nonce)),
		[]byte(fmt.Sprintf("%x", block.Timestamp)),
		block.MerkleRoot,
		block.PreviousHash,
		TransactionsToByte(transactions),
	}, []byte{})
	hash := sha256.Sum256(info)

	return hash[:]
}

func Hash(block *Models.Block) []byte {
	info := bytes.Join([][]byte{
		[]byte(fmt.Sprintf("%x", block.Nonce)),
		[]byte(fmt.Sprintf("%x", block.Timestamp)),
		block.MerkleRoot,
		block.PreviousHash,
		TransactionsToByte(block.Transactions),
	}, []byte{})
	hash := sha256.Sum256(info)

	return hash[:]
}

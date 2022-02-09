package components

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"
)

type BlockType struct {
	Nonce        int
	Timestamp    int64
	PreviousHash []byte
	CurrentHash  []byte
	Height       int
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

func CreateBlock(previousHash []byte, height int) *BlockType {
	block := &BlockType{0, time.Now().Unix(), previousHash, []byte{}, height}
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

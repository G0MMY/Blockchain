package Models

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"math/big"
	"time"
)

var (
	difficulty = 14
)

type Block struct {
	Index        int
	Nonce        int
	Timestamp    int64
	MerkleRoot   *Tree
	PreviousHash []byte
	Transactions []*Transaction
}

func CreateGenesisBlock(address []byte) *Block {
	coinbase := CreateCoinbase(address)

	block := &Block{0, 0, time.Now().Unix(), CreateTree([]*Transaction{coinbase}), []byte{}, []*Transaction{coinbase}}
	block.Proof()
	block.LinkCoinbase()

	return block
}

func (block *Block) LinkCoinbase() {
	for _, transaction := range block.Transactions {
		if transaction.IsCoinbase() {
			hash := block.Hash()
			transaction.Inputs[0].OutputTransactionId = hash
			transaction.Inputs[0].OutputBlockId = hash
		}
	}
}

func (block *Block) Proof() {
	var intHash big.Int
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	for true {
		block.Nonce += 1
		hash := block.Hash()
		intHash.SetBytes(hash)
		if intHash.Cmp(target) == -1 {
			break
		}
	}
}

func (block *Block) Hash() []byte {
	hash := sha256.Sum256(block.EncodeBlock())

	return hash[:]
}

func DecodeBlock(byteBlock []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(byteBlock))

	if err := decoder.Decode(&block); err != nil {
		log.Panic(err)
	}

	return &block
}

func (block *Block) EncodeBlock() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	if err := encoder.Encode(block); err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}

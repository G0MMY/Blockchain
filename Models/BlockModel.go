package Models

import (
	"bytes"
	"crypto/sha256"
	"github.com/ugorji/go/codec"
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
	MerkleRoot   []byte
	PreviousHash []byte
	Transactions []*Transaction
	MerkleTree   *Tree
}

func CreateGenesisBlock(privateKey []byte) *Block {
	if !IsValidPrivateKey(privateKey) {
		log.Println("Invalid private key")
		return nil
	}
	coinbase := CreateCoinbase(privateKey)

	block := &Block{0, 0, time.Now().Unix(), nil, []byte{}, []*Transaction{coinbase}, nil}
	block.addTree()
	block.proof()
	if block.Nonce == -1 {
		return nil
	}

	return block
}

func CreateBlock(privateKey []byte, index int, lastHash []byte, transactions []*Transaction) *Block {
	transactions = FindBestMemPoolTransactions(transactions, NumberOfTransactions, privateKey)
	if transactions == nil {
		return nil
	}
	coinbase := CreateCoinbase(privateKey)

	block := &Block{index, 0, time.Now().Unix(), nil, lastHash, append(transactions, coinbase), nil}
	block.addTree()
	if block.MerkleRoot == nil {
		return nil
	}

	block.proof()
	if block.Nonce == -1 {
		return nil
	}

	return block
}

func (block *Block) Validate() bool {
	if !block.ValidateProof() {
		return false
	} else if bytes.Compare(block.MerkleRoot, block.MerkleTree.RootNode.Data) != 0 {
		return false
	}

	for _, transaction := range block.Transactions {
		if !transaction.ValidateTransaction(false) {
			return false
		}
	}

	if !block.MerkleTree.CheckTree(block.Transactions) {
		return false
	}

	return true
}

func (block *Block) addTree() {
	tree := CreateTree(block.Transactions)
	if tree == nil {
		log.Println("No transactions in the block")
		return
	}

	block.MerkleTree = tree
	block.MerkleRoot = tree.RootNode.Data
}

func (block *Block) ValidateProof() bool {
	var intHash big.Int
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	hash := block.Hash()
	if hash == nil {
		return false
	}

	intHash.SetBytes(hash)
	if intHash.Cmp(target) == -1 {
		return true
	}

	return false
}

func (block *Block) proof() {
	var intHash big.Int
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	for true {
		block.Nonce += 1
		hash := block.Hash()
		if hash == nil {
			block.Nonce = -1
			return
		}

		intHash.SetBytes(hash)
		if intHash.Cmp(target) == -1 {
			break
		}
	}
}

func (block *Block) Hash() []byte {
	byteBlock := block.EncodeBlock()
	if byteBlock == nil {
		return nil
	}

	hash := sha256.Sum256(byteBlock)

	return hash[:]
}

func (block *Block) HashTransactions() [][]byte {
	var hashTransactions [][]byte

	for _, transaction := range block.Transactions {
		hashTransactions = append(hashTransactions, transaction.GetMemPoolHash(block))
	}

	return hashTransactions
}

func DecodeBlock(byteBlock []byte) *Block {
	var block BlockRequest
	decoder := codec.NewDecoder(bytes.NewReader(byteBlock), new(codec.JsonHandle))

	if err := decoder.Decode(&block); err != nil {
		log.Println(err)
		return nil
	}

	return block.CreateBlock()
}

func (block *Block) EncodeBlock() []byte {
	var buffer bytes.Buffer
	encoder := codec.NewEncoder(&buffer, new(codec.JsonHandle))

	if err := encoder.Encode(block.CreateBlockRequest()); err != nil {
		log.Println(err)
		return nil
	}

	return buffer.Bytes()
}

func (block *Block) CheckMerkleRoot(merkleRoot []byte) bool {
	if bytes.Compare(block.MerkleRoot, merkleRoot) == 0 {
		return true
	}

	return false
}

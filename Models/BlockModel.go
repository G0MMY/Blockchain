package Models

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
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

func CreateBlockToBlock(createBlock CreateBlockRequest) *Block {
	var transactions []*Transaction

	for _, stringHash := range createBlock.Transactions {
		hash, err := hex.DecodeString(stringHash)
		if err != nil {
			return nil
		}
		transactions = append(transactions, DecodeTransaction(hash))
	}

	merkleRoot, err := hex.DecodeString(createBlock.MerkleRoot)
	if err != nil {
		return nil
	}

	previousHash, err := hex.DecodeString(createBlock.PreviousHash)
	if err != nil {
		return nil
	}

	tree, err := hex.DecodeString(createBlock.MerkleTree)
	if err != nil {
		return nil
	}

	return &Block{createBlock.Index, createBlock.Nonce, createBlock.Timestamp, merkleRoot, previousHash, transactions, DecodeTree(tree)}
}

func CreateGenesisBlock(privateKey []byte) *Block {
	if !IsValidPrivateKey(privateKey) {
		log.Panic("Invalid private key")
	}
	coinbase := CreateCoinbase(privateKey)

	block := &Block{0, 0, time.Now().Unix(), nil, []byte{}, []*Transaction{coinbase}, nil}
	block.linkCoinbase()
	block.addTree()
	block.Proof()

	return block
}

func CreateBlock(privateKey []byte, index int, lastHash []byte, transactions []*Transaction) *Block {
	transactions, _ = FindBestMemPoolTransactions(transactions, NumberOfTransactions, privateKey)
	coinbase := CreateCoinbase(privateKey)

	block := &Block{index, 0, time.Now().Unix(), nil, lastHash, append(transactions, coinbase), nil}
	block.linkCoinbase()
	block.linkOutputs()
	block.addTree()
	block.Proof()

	return block
}

func (block *Block) Validate() bool {
	if !block.ValidateProof() {
		return false
	} else if bytes.Compare(block.MerkleRoot, block.MerkleTree.RootNode.Data) != 0 {
		return false
	}

	for _, transaction := range block.Transactions {
		transaction.ValidateTransaction()
	}

	if !block.MerkleTree.CheckTree(block.Transactions) {
		return false
	}

	return true
}

func (block *Block) addTree() {
	tree := CreateTree(block.Transactions)
	if tree == nil {
		log.Panic("No transactions in the block")
	}

	block.MerkleTree = tree
	block.MerkleRoot = tree.RootNode.Data
}

func (block *Block) linkOutputs() {
	for i, transaction := range block.Transactions {
		for j, output := range transaction.Outputs {
			output.Index = j
			output.TransactionIndex = i
			output.BlockId = block.Hash()
		}
	}
}

func (block *Block) linkCoinbase() {
	for i, transaction := range block.Transactions {
		if transaction.IsCoinbase() {
			transaction.Outputs[0].TransactionIndex = i
			transaction.Outputs[0].BlockId = block.Hash()
		}
	}
}

func (block *Block) ValidateProof() bool {
	var intHash big.Int
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	intHash.SetBytes(block.Hash())
	if intHash.Cmp(target) == -1 {
		return true
	}

	return false
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

func (block *Block) CheckMerkleRoot(merkleRoot []byte) bool {
	if bytes.Compare(block.MerkleRoot, merkleRoot) == 0 {
		return true
	}

	return false
}

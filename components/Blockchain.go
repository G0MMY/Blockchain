package components

import (
	"fmt"
)

type BlockchainType struct {
	Chain   []*BlockType
	Length  int
	MemPool MemPool
}

type Blockchain interface {
	AddBlock()
	AddGenesisBlock()
	DisplayBlockchain()
	GetLength() int
	GetChain() []*BlockType
	IsChainValid() bool
	AddTransaction(string, string, int, int)
}

//func InitializeBlockchain() *BlockchainType {
//	var chain []*BlockType
//	var transactions []*TransactionType
//	var memPool = &MemPoolType{transactions, "01", 0}
//	blockchain := &BlockchainType{Chain: chain, Length: 0, MemPool: memPool}
//	blockchain.AddGenesisBlock()
//
//	return blockchain
//}

func (blockchain *BlockchainType) IsChainValid() bool {
	i := 0
	for i < blockchain.Length-1 {
		if !blockchain.Chain[i].CheckBlock() {
			return false
		} else if fmt.Sprintf("%x", blockchain.Chain[i].CurrentHash) != fmt.Sprintf("%x", blockchain.Chain[i+1].PreviousHash) {
			return false
		}
		i += 1
	}
	return true
}

func (blockchain *BlockchainType) AddTransaction(sender string, receiver string, amount int, fee int) {
	blockchain.MemPool.addTransaction(CreateTransaction(sender, receiver, amount, fee))
}

func (blockchain *BlockchainType) GetChain() []*BlockType {
	return blockchain.Chain
}

func (blockchain *BlockchainType) GetLength() int {
	return blockchain.Length
}

//func (blockchain *BlockchainType) AddGenesisBlock() {
//	if blockchain.Length == 0 {
//		block := CreateBlock([]byte{0}, blockchain.MemPool.getTransactions(), 0)
//		blockchain.Chain = append(blockchain.Chain, block)
//		blockchain.Length += 1
//	}
//}
//
//func (blockchain *BlockchainType) AddBlock() {
//	height := 5
//	block := CreateBlock(blockchain.Chain[blockchain.Length-1].CurrentHash, blockchain.MemPool.getTransactions(), height)
//	blockchain.MemPool.deleteNFirstTransactions(height)
//	blockchain.Chain = append(blockchain.Chain, block)
//	blockchain.Length += 1
//}

func (blockchain *BlockchainType) DisplayBlockchain() {
	for _, block := range blockchain.Chain {
		fmt.Printf("{ \n nonce: %d, \n", block.Nonce)
		fmt.Printf(" timestamp: %d, \n", block.Timestamp)
		fmt.Printf(" previousHash: %x, \n", block.PreviousHash)
		fmt.Printf(" currentHash: %x, \n }, \n", block.CurrentHash)
	}
}

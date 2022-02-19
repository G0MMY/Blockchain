package components

import (
	"fmt"
)

type Blockchain struct {
	Chain  []*Block
	Length int
	//MemPool MemPool
}

//func InitializeBlockchain() *Blockchain {
//	var chain []*Block
//	var transactions []*TransactionType
//	var memPool = &MemPoolType{transactions, "01", 0}
//	blockchain := &Blockchain{Chain: chain, Length: 0, MemPool: memPool}
//	blockchain.AddGenesisBlock()
//
//	return blockchain
//}

//func (blockchain *Blockchain) AddTransaction(sender string, receiver string, amount int, fee int) {
//	blockchain.MemPool.addTransaction(CreateTransaction(sender, receiver, amount, fee))
//}

//func (blockchain *Blockchain) AddGenesisBlock() {
//	if blockchain.Length == 0 {
//		block := CreateBlock([]byte{0}, blockchain.MemPool.getTransactions(), 0)
//		blockchain.Chain = append(blockchain.Chain, block)
//		blockchain.Length += 1
//	}
//}
//
//func (blockchain *Blockchain) AddBlock() {
//	height := 5
//	block := CreateBlock(blockchain.Chain[blockchain.Length-1].CurrentHash, blockchain.MemPool.getTransactions(), height)
//	blockchain.MemPool.deleteNFirstTransactions(height)
//	blockchain.Chain = append(blockchain.Chain, block)
//	blockchain.Length += 1
//}

func (blockchain *Blockchain) DisplayBlockchain() {
	for _, block := range blockchain.Chain {
		fmt.Printf("{ \n nonce: %d, \n", block.Nonce)
		fmt.Printf(" timestamp: %d, \n", block.Timestamp)
		fmt.Printf(" previousHash: %x, \n", block.PreviousHash)
		fmt.Printf(" currentHash: %x, \n }, \n", block.CurrentHash)
	}
}

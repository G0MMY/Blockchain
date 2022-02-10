package main

import (
	"blockchain/components"
)

func main() {
	var chain []*components.BlockType
	var blockchain components.Blockchain
	var transactions []*components.TransactionType
	var memPool = &components.MemPoolType{transactions, "01", 0}
	blockchain = &components.BlockchainType{Chain: chain, Length: 0, MemPool: memPool}
	blockchain.AddGenesisBlock()
	blockchain.AddBlock()

	blockchain.AddTransaction("me", "me", 150, 10)
	blockchain.AddTransaction("me", "me", 1096, 100)
	blockchain.AddTransaction("me", "me", 12, 5)
	blockchain.AddTransaction("me", "me", 10, 1)
	blockchain.AddTransaction("me", "me", 1086, 180)
	blockchain.AddTransaction("me", "me", 45, 4)

	blockchain.AddBlock()

	blockchain.AddTransaction("me", "me", 100, 10)
	blockchain.AddTransaction("me", "me", 106, 1)
	blockchain.AddTransaction("me", "me", 102, 50)

	blockchain.AddBlock()
	blockchain.AddBlock()
}

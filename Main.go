package main

import (
	"blockchain/components"
)

func main() {
	var chain []*components.BlockType
	var blockchain components.Blockchain
	blockchain = &components.BlockchainType{Chain: chain, Length: 0}
	blockchain.AddGenesisBlock()
	blockchain.AddBlock()
	blockchain.AddBlock()
	blockchain.AddBlock()
	blockchain.AddBlock()
	blockchain.DisplayBlockchain()
}

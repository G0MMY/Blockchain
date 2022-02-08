package main

import (
	"blockchain/components"
)

func main() {
	var chain []components.Block
	var blockchain components.Blockchain = components.Blockchain{Chain: chain}
	blockchain.Chain = components.AddBlock(blockchain)
	blockchain.Chain = components.AddBlock(blockchain)
	blockchain.Chain = components.AddBlock(blockchain)
	blockchain.Chain = components.AddBlock(blockchain)
	components.DisplayBlockchain(blockchain)
}

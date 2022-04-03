package main

import (
	"blockchain/Models"
	"fmt"
)

func main() {
	blockchain := Models.InitBlockchain([]byte("test"))

	//blockchain.CreateBlock([]byte("test"))
	//blockchain.CreateBlock([]byte("test"))
	//blockchain.CreateBlock([]byte("test"))
	//blockchain.CreateBlock([]byte("test"))
	//blockchain.CreateBlock([]byte("test"))

	chain := blockchain.GetBlockchain()

	for _, block := range chain {
		fmt.Printf("block height: %d\n", block.Index)
		fmt.Printf("block hash: %x\n", block.Hash())
		fmt.Println()
		fmt.Println()
	}

	blockchain.DB.Close()
}

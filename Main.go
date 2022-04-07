package main

import (
	"blockchain/Models"
	"fmt"
	"time"
)

func main() {
	blockchain := Models.InitBlockchain([]byte("max"))

	blockchain.CreateTransaction([]byte("max"), []byte("joe"), 5, 5, time.Now().Unix())
	blockchain.CreateTransaction([]byte("max"), []byte("joe"), 5, 5, time.Now().Unix())
	blockchain.CreateTransaction([]byte("max"), []byte("joe"), 5, 5, time.Now().Unix())

	blockchain.CreateBlock([]byte("joe"))

	chain := blockchain.GetBlockchain()

	for _, block := range chain {
		fmt.Printf("block height: %d\n", block.Index)
		fmt.Printf("block hash: %x\n", block.Hash())
		fmt.Println()
	}

	blockchain.DB.Close()
}

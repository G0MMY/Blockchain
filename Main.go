package main

import (
	"blockchain/Models"
	"fmt"
)

func main() {
	blockchain := Models.InitBlockchain([]byte("test"))

	blockchain.CreateBlock()
	blockchain.CreateBlock()
	blockchain.CreateBlock()
	blockchain.CreateBlock()

	chain := blockchain.GetBlockchain()

	fmt.Println(chain)
}

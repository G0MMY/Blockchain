package main

import (
	"blockchain/Models"
	"blockchain/Network"
)

func main() {
	//cmd.Execute()
	//wallet1 := Models.CreateWallet()
	//
	//blockchain := Models.InitBlockchain(wallet1.PrivateKey)
	//
	//chain := blockchain.GetBlockchain()
	//
	//for _, block := range chain {
	//	fmt.Printf("block height: %d\n", block.Index)
	//	fmt.Printf("block hash: %x\n", block.Hash())
	//	fmt.Println()
	//}
	//
	//blockchain.DB.Close()
	Network.InitializeNode(Models.CreateWallet().PrivateKey)
}

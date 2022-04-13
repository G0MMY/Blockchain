/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"blockchain/Models"
)

func main() {
	//cmd.Execute()
	wallet1 := Models.CreateWallet()
	//wallet2 := Models.CreateWallet()

	blockchain := Models.InitBlockchain(wallet1.PrivateKey)
	//blockchain.CreateBlock(wallet1.PrivateKey)
	//blockchain.CreateBlock(wallet1.PrivateKey)
	//blockchain.CreateBlock(wallet1.PrivateKey)
	//blockchain.CreateBlock(wallet1.PrivateKey)
	//
	//blockchain.CreateTransaction(wallet1.PrivateKey, wallet2.PublicKey, 10, 10, time.Now().Unix())
	//blockchain.CreateTransaction(wallet1.PrivateKey, wallet2.PublicKey, 10, 10, time.Now().Unix())
	//blockchain.CreateTransaction(wallet1.PrivateKey, wallet2.PublicKey, 10, 10, time.Now().Unix())
	//blockchain.CreateTransaction(wallet1.PrivateKey, wallet2.PublicKey, 10, 10, time.Now().Unix())
	//blockchain.CreateTransaction(wallet1.PrivateKey, wallet2.PublicKey, 10, 10, time.Now().Unix())

	t := blockchain.GetMemPoolTransactions()

	Models.CreateTree(t)

	//blockchain.CreateBlock(wallet1.PrivateKey)
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
}

//export PATH=$PATH:/home/maxim/go/bin

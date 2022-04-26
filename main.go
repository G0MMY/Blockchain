package main

import "blockchain/Network"

func main() {
	//cmd.Execute()
	//wallet1 := Models.CreateWallet()
	//wallet2 := Models.CreateWallet()
	//
	//blockchain := Models.InitBlockchain("4000")
	//
	//lastBlock := blockchain.GetLastBlock()
	//memPool := blockchain.GetMemPoolTransactions()
	//b := Models.CreateBlock(wallet1.PrivateKey, lastBlock.Index+1, lastBlock.Hash(), memPool)
	//blockchain.AddBlock(b)
	//blockchain.CreateTransaction(wallet1.PrivateKey, wallet2.PublicKey, 10, 10, time.Now().Unix())
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
	//priv, _ := hex.DecodeString("307702010104207d8dc86cd1956e3a4a57aaa3681ccfd9422ca4420046f9ea89417e041372d5d0a00a06082a8648ce3d030107a14403420004b0fb2b18f23b976e47938c7b9b3c232838f43eb42e916445a754a897ee8c98c0b3eb73b39affc87f96dae1f92659f23a031d411f3b2bbef09bd2288f9dccd9c4")

	//Network.InitializeNode("4000", []string{})
	//Network.InitializeNode("4001", []string{"4000"})
	Network.InitializeNode("4002", []string{"4000", "4001"})
}

//3059301306072a8648ce3d020106082a8648ce3d03010703420004b0fb2b18f23b976e47938c7b9b3c232838f43eb42e916445a754a897ee8c98c0b3eb73b39affc87f96dae1f92659f23a031d411f3b2bbef09bd2288f9dccd9c4
//313274356b4251457445707257556f593648774d6270487876565233726d4c43746a

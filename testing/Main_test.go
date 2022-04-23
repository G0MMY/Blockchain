package testing

import (
	"blockchain/Models"
	"bytes"
	"testing"
	"time"
)

var (
	wallet1    *Models.Wallet
	wallet2    *Models.Wallet
	blockchain *Models.Blockchain
)

func TestInit(t *testing.T) {
	wallet1 = Models.CreateWallet()
	wallet2 = Models.CreateWallet()
	blockchain = Models.InitTestBlockchain(wallet1.PrivateKey)
}

func TestInitBlockchain(t *testing.T) {
	if blockchain == nil {
		t.Error("blockchain in not initialized")
	} else if blockchain.DB == nil {
		t.Error("The blockchain database is nil")
	} else if blockchain.LastHash == nil {
		t.Error("The blockchain lastHash is nil")
	}
}

func TestCreateBlock(t *testing.T) {
	if blockchain == nil {
		t.Error("blockchain in not initialized")
		return
	}
	lastBlock := blockchain.GetLastBlock()
	block := Models.CreateBlock(wallet1.PrivateKey, lastBlock.Index+1, lastBlock.Hash(), blockchain.GetMemPoolTransactions())
	block, err := blockchain.CreateBlock(wallet1.PrivateKey, block)

	if err != "" {
		t.Error(err)
	} else if bytes.Compare(blockchain.LastHash, block.Hash()) != 0 {
		t.Error("The blockchain lastHash hasn't been updated properly")
	} else if bytes.Compare(lastBlock.Hash(), block.PreviousHash) != 0 {
		t.Error("The new block doesn't have the right previous hash")
	} else if len(block.Transactions) != 1 {
		t.Errorf("The block as %d transactions instead of 1", len(block.Transactions))
	} else if len(block.Transactions[0].Outputs) != 1 {
		t.Errorf("The block as %d transaction outputs instead of 1", len(block.Transactions[0].Outputs))
	} else if len(block.Transactions[0].Inputs) != 1 {
		t.Errorf("The block as %d transaction inputs instead of 1", len(block.Transactions[0].Inputs))
	} else if block.MerkleRoot == nil {
		t.Error("The block has no merkle root")
	} else if !block.MerkleTree.CheckTree(block.Transactions) {
		t.Error("Merkle tree or transactions are invalid")
	} else if bytes.Compare(block.MerkleRoot, block.MerkleTree.RootNode.Data) != 0 {
		t.Error("The merkle tree root is not the same as the block merkle root")
	} else if !block.ValidateProof() {
		t.Error("The block proof isn't good")
	}
}

func TestCreateTransaction(t *testing.T) {
	if blockchain == nil {
		t.Error("blockchain in not initialized")
		return
	}
	unspentOutputs := blockchain.GetUnspentOutputs(Models.ValidateAddress(wallet1.PublicKey))
	transaction := blockchain.CreateTransaction(wallet1.PrivateKey, wallet2.PublicKey, 1, 5, time.Now().Unix())
	unspentOutputsAfter := blockchain.GetUnspentOutputs(wallet1.PublicKey)

	if unspentOutputs != nil && unspentOutputsAfter != nil {
		if len(unspentOutputs.Outputs) <= len(unspentOutputsAfter.Outputs) {
			t.Error("The unspent outputs didn't change")
			return
		}
	}
	for _, input := range transaction.Inputs {
		if input.Output == nil {
			t.Errorf("The input isn't linked to an output")
		} else if bytes.Compare(input.Output.PublicKeyHash, Models.ValidateAddress(wallet1.PublicKey)) != 0 {
			t.Errorf("There is an input that dosen't have the right publicKey hash")
		} else if !Models.ValidateSignature(input.Output.Amount, input.PublicKey, input.Signature) {
			t.Error("The transaction isn't signed properly")
		}
	}
}

func TestTransactionsInBlock(t *testing.T) {
	if blockchain == nil {
		t.Error("blockchain in not initialized")
		return
	}
	memPool := blockchain.GetMemPoolTransactions()

	if memPool == nil || len(memPool) == 0 {
		t.Error("There are no transactions")
		return
	}

	lastBlock := blockchain.GetLastBlock()
	block := Models.CreateBlock(wallet1.PrivateKey, lastBlock.Index+1, lastBlock.Hash(), memPool)
	block, err := blockchain.CreateBlock(wallet1.PrivateKey, block)

	if err != "" {
		t.Error(err)
		return
	} else if len(block.Transactions) != len(memPool)+1 {
		t.Error("Transaction missing")
		return
	}

	hasCoinbase := false
	fees := 0
	for _, transaction := range block.Transactions {
		transaction.ValidateTransaction()
		if !hasCoinbase && transaction.IsCoinbase() {
			hasCoinbase = true
		}
		for _, output := range transaction.Outputs {
			if output.Amount == transaction.Fee && bytes.Compare(output.PublicKeyHash, Models.GetPublicKeyHash(wallet1.PublicKey)) == 0 {
				fees += 1
				break
			}
		}
	}

	if !hasCoinbase {
		t.Error("The block as no coinbase transaction")
	}
	if fees != len(block.Transactions)-1 {
		t.Error("The block contains one or more transactions with no fee to miner")
	}
}

func TestMultipleTransactions(t *testing.T) {
	blockchain.DB.Close()
	blockchain = Models.InitTestBlockchain(wallet1.PrivateKey)

	i := 0
	for i < Models.NumberOfTransactions*2 {
		lastBlock := blockchain.GetLastBlock()
		memPool := blockchain.GetMemPoolTransactions()
		block := Models.CreateBlock(wallet1.PrivateKey, lastBlock.Index+1, lastBlock.Hash(), memPool)
		block, err := blockchain.CreateBlock(wallet1.PrivateKey, block)

		if err != "" {
			t.Error(err)
			return
		}
		i += 1
	}
	i = 0
	for i < Models.NumberOfTransactions*2 {
		blockchain.CreateTransaction(wallet1.PrivateKey, wallet2.PublicKey, 30, i, time.Now().Unix())
		i += 1
	}
	lastBlock := blockchain.GetLastBlock()
	memPool := blockchain.GetMemPoolTransactions()
	block1 := Models.CreateBlock(wallet2.PrivateKey, lastBlock.Index+1, lastBlock.Hash(), memPool)
	block1, err := blockchain.CreateBlock(wallet2.PrivateKey, block1)

	if err != "" {
		t.Error(err)
		return
	}

	lastBlock = blockchain.GetLastBlock()
	memPool = blockchain.GetMemPoolTransactions()
	block2 := Models.CreateBlock(wallet2.PrivateKey, lastBlock.Index+1, lastBlock.Hash(), memPool)
	block2, err = blockchain.CreateBlock(wallet2.PrivateKey, block2)

	if err != "" {
		t.Error(err)
	} else if len(block1.Transactions) != Models.NumberOfTransactions {
		t.Errorf("The block 1 has %d transactions but needed %d", len(block1.Transactions), Models.NumberOfTransactions)
	} else if len(block2.Transactions) != Models.NumberOfTransactions {
		t.Errorf("The block 2 has %d transactions but needed %d", len(block2.Transactions), Models.NumberOfTransactions)
	} else if !block1.MerkleTree.CheckTree(block1.Transactions) {
		t.Error("Invalid merkle tree or transactions in block 1")
	} else if !block2.MerkleTree.CheckTree(block2.Transactions) {
		t.Error("Invalid merkle tree or transactions in block 2")
	} else {
		target1 := 0
		feeBlock1 := 0
		for m, transaction := range block1.Transactions {
			if m != len(block1.Transactions)-1 {
				target1 += Models.NumberOfTransactions*2 - m - 1
			}
			feeBlock1 += transaction.Fee
		}

		target2 := 0
		feeBlock2 := 0
		for l, transaction := range block2.Transactions {
			if l != len(block1.Transactions)-1 {
				target2 += Models.NumberOfTransactions - l
			}
			feeBlock2 += transaction.Fee
		}

		if feeBlock1 != target1 {
			t.Error("The block 1 doesn't have the right transactions")
		} else if feeBlock2 != target2 {
			t.Error("The block 2 doesn't have the right transactions")
		}
	}
}

func TestEnd(t *testing.T) {
	blockchain.DB.Close()
}

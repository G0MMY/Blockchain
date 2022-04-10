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
	blockchain = Models.InitTestBlockchain(wallet1.PublicKey)
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

//add checks
func TestCreateBlock(t *testing.T) {
	if blockchain == nil {
		t.Error("blockchain in not initialized")
	}
	lastHash := blockchain.LastHash
	block := blockchain.CreateBlock(wallet1.PublicKey)

	if bytes.Compare(blockchain.LastHash, block.Hash()) != 0 {
		t.Error("The blockchain lastHash asn't been updated properly")
	} else if bytes.Compare(lastHash, block.PreviousHash) != 0 {
		t.Error("The new block dosen't have the right previous hash")
	} else if len(block.Transactions) != 1 {
		t.Errorf("The block as %d transactions instead of 1", len(block.Transactions))
	} else if len(block.Transactions[0].Outputs) != 1 {
		t.Errorf("The block as %d transaction outputs instead of 1", len(block.Transactions[0].Outputs))
	} else if len(block.Transactions[0].Inputs) != 1 {
		t.Errorf("The block as %d transaction inputs instead of 1", len(block.Transactions[0].Inputs))
	}
}

func TestCreateTransaction(t *testing.T) {
	if blockchain == nil {
		t.Error("blockchain in not initialized")
	}
	unspentOutputs := blockchain.GetUnspentOutputs(Models.ValidateAddress(wallet1.PublicKey))
	t.Errorf("%d", len(unspentOutputs.Outputs))
	transaction1 := blockchain.CreateTransaction(wallet1.PublicKey, wallet2.PublicKey, 1, 5, time.Now().Unix())
	pubKeyHash := Models.GetPublicKeyHash(wallet1.PublicKey)

	for _, output := range transaction1.Outputs {
		if bytes.Compare(output.PublicKeyHash, pubKeyHash) != 0 {
			t.Errorf("The publicKeyHash in output %d is invalid", output.Index)
		}
	}

	unspentOutputsAfter := blockchain.GetUnspentOutputs(wallet1.PublicKey)

	if len(unspentOutputs.Outputs) <= len(unspentOutputsAfter.Outputs) {
		t.Error("The unspent outputs didin't change")
	}
}

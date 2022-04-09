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
	lastHash := blockchain.LastHash
	block := blockchain.CreateBlock(wallet1.PublicKey)

	if bytes.Compare(blockchain.LastHash, block.Hash()) != 0 {
		t.Error("The blockchain lastHash asn't been updated properly")
	} else if bytes.Compare(lastHash, block.PreviousHash) != 0 {
		t.Error("The new block dosen't have the right previous hash")
	}
}

func TestCreateTransaction(t *testing.T) {
	transaction1 := blockchain.CreateTransaction(wallet1.PublicKey, wallet2.PublicKey, 1, 5, time.Now().Unix())
	transaction2 := blockchain.CreateTransaction(wallet1.PublicKey, wallet2.PublicKey, 2, 5, time.Now().Unix())
	transaction3 := blockchain.CreateTransaction(wallet1.PublicKey, wallet2.PublicKey, 3, 5, time.Now().Unix())

}

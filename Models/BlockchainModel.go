package Models

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"log"
	"os"
)

var (
	NumberOfTransactions = 5
)

type Blockchain struct {
	LastHash []byte
	DB       *leveldb.DB
}

type BlockchainIterator struct {
	CurrentHash []byte
	DB          *leveldb.DB
}

func InitTestBlockchain(privateKey []byte) *Blockchain {
	if !IsValidPrivateKey(privateKey) {
		log.Panic("Invalid private key")
	}
	var read *opt.ReadOptions

	db, err := leveldb.OpenFile("./testing/dbTest", nil)
	if err != nil {
		log.Panic(err)
	}

	hasLastHash, err := db.Has([]byte("lastHash"), read)
	if err != nil {
		log.Panic(err)
	}

	if hasLastHash {
		err = os.RemoveAll("./testing/dbTest")
		if err != nil {
			log.Panic(err)
		}

		db, err = leveldb.OpenFile("./testing/dbTest", nil)
		if err != nil {
			log.Panic(err)
		}
	}

	block := CreateGenesisBlock(privateKey)
	blockchain := &Blockchain{[]byte{}, db}
	blockchain.addBlock(block)

	return blockchain
}

func InitBlockchain(privateKey []byte) *Blockchain {
	if !IsValidPrivateKey(privateKey) {
		log.Panic("Invalid private key")
	}
	var read *opt.ReadOptions
	db, err := leveldb.OpenFile("./db", nil)

	if err != nil {
		log.Panic(err)
	}

	hasLastHash, err := db.Has([]byte("lastHash"), read)

	if err != nil {
		log.Panic(err)
	}

	if hasLastHash {
		lastHash, err := db.Get([]byte("lastHash"), read)

		if err != nil {
			log.Panic(err)
		}

		return &Blockchain{lastHash, db}
	}

	block := CreateGenesisBlock(privateKey)
	blockchain := &Blockchain{[]byte{}, db}
	blockchain.addBlock(block)

	return blockchain
}

func (blockchain *Blockchain) GetBlockchain() []*Block {
	blockchainIterator := &BlockchainIterator{blockchain.LastHash, blockchain.DB}

	return blockchainIterator.getBlockchain()
}

func (iter *BlockchainIterator) getBlockchain() []*Block {
	var blockchain []*Block

	currentBlock := iter.next()
	if currentBlock != nil {
		blockchain = append(blockchain, currentBlock)

		for bytes.Compare(iter.CurrentHash, []byte{}) != 0 {
			blockchain = append(blockchain, iter.next())
		}
	}

	return blockchain
}

func (iter *BlockchainIterator) next() *Block {
	var read *opt.ReadOptions

	if byteBlock, err := iter.DB.Get(iter.CurrentHash, read); err != nil {
		log.Panic(err)
	} else {
		currentBlock := DecodeBlock(byteBlock)

		if bytes.Compare(currentBlock.Hash(), iter.CurrentHash) != 0 {
			log.Panic("The chain is invalid")
		}
		iter.CurrentHash = currentBlock.PreviousHash

		return currentBlock
	}

	return nil
}

func (blockchain *Blockchain) GetLastBlock() *Block {
	if blockchain.DB != nil {
		var read *opt.ReadOptions

		if blockByte, err := blockchain.DB.Get(blockchain.LastHash, read); err != nil {
			log.Panic(err)
		} else {
			return DecodeBlock(blockByte)
		}
	}
	return nil
}

func (blockchain *Blockchain) HashMemPoolTransactions() []byte {
	var hashTransactions [][]byte

	for _, transaction := range blockchain.GetMemPoolTransactions() {
		hashTransactions = append(hashTransactions, transaction.Hash())
	}

	hash := sha256.Sum256(bytes.Join(hashTransactions, []byte{}))

	return hash[:]
}

func (blockchain *Blockchain) GetBlock(blockHash []byte) *Block {
	var read *opt.ReadOptions

	if blockByte, err := blockchain.DB.Get(blockHash, read); err != nil {
		if err.Error() == "leveldb: not found" {
			return nil
		}
		log.Panic(err)
	} else {
		return DecodeBlock(blockByte)
	}

	return nil
}

func (blockchain *Blockchain) GetMerkleRoot(blockHash []byte) []byte {
	return blockchain.GetBlock(blockHash).MerkleRoot
}

func (blockchain *Blockchain) CreateBlock(privateKey []byte) *Block {
	if !IsValidPrivateKey(privateKey) {
		log.Panic("Invalid private key")
	}
	lastBlock := blockchain.GetLastBlock()
	transactions, transactionsHash := FindBestMemPoolTransactions(blockchain.GetMemPoolTransactions(), NumberOfTransactions, privateKey)

	blockchain.updateMemPoolTransactions(transactionsHash)

	if lastBlock != nil {
		block := CreateBlock(privateKey, lastBlock.Index+1, blockchain.LastHash, transactions)

		blockchain.addBlock(block)
		return block
	}
	return nil
}

func (blockchain *Blockchain) addBlock(block *Block) {
	var write *opt.WriteOptions
	hash := block.Hash()

	blockchain.persistUnspentOutputs(block)

	if err := blockchain.DB.Put([]byte("lastHash"), hash, write); err != nil {
		log.Panic(err)
	}

	if err := blockchain.DB.Put(hash, block.EncodeBlock(), write); err != nil {
		log.Panic(err)
	}

	blockchain.LastHash = hash
}

func (blockchain *Blockchain) persistUnspentOutputs(block *Block) {
	var publicKeys []string
	outputsPerPubKey := make(map[string][]*Output)

	for _, transaction := range block.Transactions {
		for _, output := range transaction.Outputs {
			publicKeyString := fmt.Sprintf("%s", output.PublicKeyHash)
			if _, ok := outputsPerPubKey[publicKeyString]; !ok {
				publicKeys = append(publicKeys, publicKeyString)
			}
			outputsPerPubKey[publicKeyString] = append(outputsPerPubKey[publicKeyString], output)
		}
	}

	for _, publicKeyString := range publicKeys {
		publicKey := []byte(publicKeyString)
		unspentOutputs := blockchain.GetUnspentOutputs(publicKey)
		if unspentOutputs == nil {
			blockchain.updateUnspentOutputs(&UnspentOutput{outputsPerPubKey[publicKeyString]}, publicKey)
		} else {
			unspentOutputs.Outputs = append(unspentOutputs.Outputs, outputsPerPubKey[publicKeyString]...)
			blockchain.updateUnspentOutputs(unspentOutputs, publicKey)
		}
	}
}

func (blockchain *Blockchain) GetUnspentOutputs(address []byte) *UnspentOutput {
	var read *opt.ReadOptions

	key := GenerateUnspentOutputKey(address)

	if outputs, err := blockchain.DB.Get(key, read); err != nil {
		if err.Error() != "leveldb: not found" {
			log.Panic()
		}
	} else {
		return DecodeUnspentOutput(outputs)
	}

	return nil
}

func (blockchain *Blockchain) updateMemPoolTransactions(memPoolTransactions [][]byte) {
	var write *opt.WriteOptions

	for _, hash := range memPoolTransactions {
		key := bytes.Join([][]byte{
			[]byte("MemPool-"),
			hash,
		}, []byte{})

		if err := blockchain.DB.Delete(key, write); err != nil {
			log.Panic(err)
		}
	}
}

func (blockchain *Blockchain) updateUnspentOutputs(unspentOuputs *UnspentOutput, address []byte) {
	var write *opt.WriteOptions

	key := GenerateUnspentOutputKey(address)

	if err := blockchain.DB.Put(key, unspentOuputs.EncodeUnspentOutput(), write); err != nil {
		log.Panic(err)
	}
}

func (blockchain *Blockchain) GetMemPoolTransactions() []*Transaction {
	var read *opt.ReadOptions
	var transactions []*Transaction

	iter := blockchain.DB.NewIterator(util.BytesPrefix([]byte("MemPool-")), read)

	for iter.Next() {
		byteTransaction := iter.Value()
		transactions = append(transactions, DecodeTransaction(byteTransaction))
	}
	iter.Release()

	return transactions
}

func (blockchain *Blockchain) CreateTransaction(privateKey, to []byte, amount, fee int, timestamp int64) *Transaction {
	var write *opt.WriteOptions

	if !IsValidPrivateKey(privateKey) {
		log.Panic("Invalid private key")
	}

	fromHash := ValidateAddress(GetPublicKeyFromPrivateKey(privateKey))
	toHash := ValidateAddress(to)

	if bytes.Compare(fromHash, toHash) == 0 {
		log.Panic("You can't send money to yourself")
	}

	unspentOutputs := blockchain.GetUnspentOutputs(fromHash)
	if unspentOutputs == nil || unspentOutputs.Outputs == nil {
		log.Panic("You have no money buddy!")
	}
	outputRest, amountRest := unspentOutputs.GetOutputsForAmount(amount + fee)

	if amountRest > 0 {
		log.Panic("You don't have enough money")
	}

	blockchain.updateUnspentOutputs(&UnspentOutput{outputRest}, fromHash)

	transaction := CreateTransaction(toHash, fromHash, privateKey, amount, amountRest, fee, timestamp, unspentOutputs)
	key := bytes.Join([][]byte{
		[]byte("MemPool-"),
		transaction.Hash(),
	}, []byte{})

	if err := blockchain.DB.Put(key, transaction.EncodeTransaction(), write); err != nil {
		log.Panic(err)
	}

	return transaction
}

func (blockchain *Blockchain) GetBalance(address []byte) int {
	pubKeyHash := ValidateAddress(address)

	return GetBalance(address, blockchain.GetUnspentOutputs(pubKeyHash))
}

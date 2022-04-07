package Models

import (
	"bytes"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"log"
)

type Blockchain struct {
	LastHash []byte
	DB       *leveldb.DB
}

type BlockchainIterator struct {
	CurrentHash []byte
	DB          *leveldb.DB
}

func InitBlockchain(address []byte) *Blockchain {
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

	block := CreateGenesisBlock(address)
	blockchain := &Blockchain{[]byte{}, db}
	blockchain.AddBlock(block)

	return blockchain
}

func (blockchain *Blockchain) GetBlockchain() []*Block {
	blockchainIterator := &BlockchainIterator{blockchain.LastHash, blockchain.DB}

	return blockchainIterator.GetBlockchain()
}

func (iter *BlockchainIterator) GetBlockchain() []*Block {
	var blockchain []*Block

	currentBlock := iter.Next()
	if currentBlock != nil {
		blockchain = append(blockchain, currentBlock)

		for bytes.Compare(iter.CurrentHash, []byte{}) != 0 {
			blockchain = append(blockchain, iter.Next())
		}
	}

	return blockchain
}

func (iter *BlockchainIterator) Next() *Block {
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

//add merkle root and transactions and transaction with fee to miner
func (blockchain *Blockchain) CreateBlock(address []byte) *Block {
	lastBlock := blockchain.GetLastBlock()
	transactions := FindBestMemPoolTransactions(blockchain.GetMemPoolTransactions(), 5)

	if lastBlock != nil {
		block := CreateBlock(address, lastBlock.Index+1, blockchain.LastHash, transactions, &Tree{})

		blockchain.AddBlock(block)
		return block
	}
	return nil
}

func (blockchain *Blockchain) AddBlock(block *Block) {
	var write *opt.WriteOptions
	hash := block.Hash()

	blockchain.PersistUnspentOutputs(block)

	if err := blockchain.DB.Put([]byte("lastHash"), hash, write); err != nil {
		log.Panic(err)
	}

	if err := blockchain.DB.Put(hash, block.EncodeBlock(), write); err != nil {
		log.Panic(err)
	}

	blockchain.LastHash = hash
}

//change from %s to %x when wallet
func (blockchain *Blockchain) PersistUnspentOutputs(block *Block) {
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
			blockchain.UpdateUnspentOutputs(&UnspentOutput{outputsPerPubKey[publicKeyString]}, publicKey)
		} else {
			unspentOutputs.Outputs = append(unspentOutputs.Outputs, outputsPerPubKey[publicKeyString]...)
			blockchain.UpdateUnspentOutputs(unspentOutputs, publicKey)
		}
	}
}

func (blockchain *Blockchain) GetUnspentOutputs(from []byte) *UnspentOutput {
	var read *opt.ReadOptions

	key := bytes.Join([][]byte{
		[]byte("UnspentOutput-"),
		from,
	}, []byte{})

	if outputs, err := blockchain.DB.Get(key, read); err != nil {
		if err.Error() != "leveldb: not found" {
			log.Panic()
		}
	} else {
		return DecodeUnspentOutput(outputs)
	}

	return nil
}

func (blockchain *Blockchain) UpdateUnspentOutputs(unspentOuputs *UnspentOutput, address []byte) {
	var write *opt.WriteOptions

	key := bytes.Join([][]byte{
		[]byte("UnspentOutput-"),
		address,
	}, []byte{})

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

func (blockchain *Blockchain) CreateTransaction(from, to []byte, amount, fee int, timestamp int64) {
	var write *opt.WriteOptions

	unspentOutputs := blockchain.GetUnspentOutputs(from)
	if unspentOutputs == nil || unspentOutputs.Outputs == nil {
		log.Panic("You have no money buddy!")
	}
	outputRest, amountRest := unspentOutputs.GetOutputsForAmount(amount + fee)
	blockchain.UpdateUnspentOutputs(&UnspentOutput{outputRest}, from)

	transaction := CreateTransaction(to, from, amount, amountRest, fee, timestamp, unspentOutputs)
	key := bytes.Join([][]byte{
		[]byte("MemPool-"),
		transaction.Hash(),
	}, []byte{})

	if err := blockchain.DB.Put(key, transaction.EncodeTransaction(), write); err != nil {
		log.Panic(err)
	}
}

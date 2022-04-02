package Models

import (
	"bytes"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"log"
	"time"
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

	res, err := db.Has([]byte("lastHash"), read)

	if err != nil {
		log.Panic(err)
	}

	if res {
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

func (blockchain Blockchain) GetBlockchain() []*Block {
	blockchainIterator := BlockchainIterator{blockchain.LastHash, blockchain.DB}

	return blockchainIterator.GetBlockchain()
}

func (iter BlockchainIterator) GetBlockchain() []*Block {
	var blockchain []*Block

	currentBlock := iter.next()
	if currentBlock != nil {
		blockchain = append(blockchain, currentBlock)

		for bytes.Compare(currentBlock.PreviousHash, []byte{}) != 0 {
			blockchain = append(blockchain, iter.next())
		}
	}

	return blockchain
}

func (iter BlockchainIterator) next() *Block {
	var read *opt.ReadOptions

	if byteBlock, err := iter.DB.Get(iter.CurrentHash, read); err != nil {
		log.Panic(err)
	} else {
		currentBlock := DecodeBlock(byteBlock)
		iter.CurrentHash = currentBlock.Hash()

		return currentBlock
	}

	return nil
}

func (blockchain Blockchain) GetLastBlock() *Block {
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

//add merkle root and transactions
func (blockchain Blockchain) CreateBlock() {
	lastBlock := blockchain.GetLastBlock()

	if lastBlock != nil {
		block := Block{lastBlock.Index + 1, 0, time.Now().Unix(), &Tree{}, blockchain.LastHash, []*Transaction{}}

		blockchain.LastHash = blockchain.AddBlock(&block)
	}
}

func (blockchain Blockchain) AddBlock(block *Block) []byte {
	var write *opt.WriteOptions
	hash := block.Hash()

	if err := blockchain.DB.Put([]byte("lastHash"), hash, write); err != nil {
		log.Panic(err)
	}

	if err := blockchain.DB.Put(hash, block.EncodeBlock(), write); err != nil {
		log.Panic(err)
	}

	return hash
}

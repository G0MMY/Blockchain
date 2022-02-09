package components

import (
	"fmt"
)

type BlockchainType struct {
	Chain  []*BlockType
	Length int
}

type Blockchain interface {
	AddBlock()
	AddGenesisBlock()
	DisplayBlockchain()
	GetLength() int
	GetChain() []*BlockType
	IsChainValid() bool
}

func (blockchain *BlockchainType) IsChainValid() bool {
	i := 0
	for i < blockchain.Length-1 {
		if !blockchain.Chain[i].CheckBlock() {
			return false
		} else if fmt.Sprintf("%x", blockchain.Chain[i].CurrentHash) != fmt.Sprintf("%x", blockchain.Chain[i+1].PreviousHash) {
			return false
		}
		i += 1
	}
	return true
}

func (blockchain *BlockchainType) GetChain() []*BlockType {
	return blockchain.Chain
}

func (blockchain *BlockchainType) GetLength() int {
	return blockchain.Length
}

func (blockchain *BlockchainType) AddGenesisBlock() {
	block := CreateBlock([]byte{0}, 0)
	blockchain.Chain = append(blockchain.Chain, block)
	blockchain.Length += 1
}

func (blockchain *BlockchainType) AddBlock() {
	block := CreateBlock(blockchain.Chain[blockchain.Length-1].CurrentHash, 10)
	blockchain.Chain = append(blockchain.Chain, block)
	blockchain.Length += 1
}

func (blockchain *BlockchainType) DisplayBlockchain() {
	for _, block := range blockchain.Chain {
		fmt.Printf("{ \n nonce: %d, \n", block.Nonce)
		fmt.Printf(" timestamp: %d, \n", block.Timestamp)
		fmt.Printf(" previousHash: %x, \n", block.PreviousHash)
		fmt.Printf(" currentHash: %x, \n }, \n", block.CurrentHash)
	}
}

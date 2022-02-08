package components

import (
	"fmt"
	"time"
)

type Blockchain struct {
	Chain []Block
}

func AddBlock(blockchain Blockchain) []Block {
	var block Block
	if len(blockchain.Chain) == 0 {
		block = MineBlock(0, 0, time.Now().Unix(), "0")
	} else {
		block = MineBlock(len(blockchain.Chain), 0, time.Now().Unix(), blockchain.Chain[len(blockchain.Chain)-1].CurrentHash)
	}
	return append(blockchain.Chain, block)
}

func DisplayBlockchain(blockchain Blockchain) {
	for _, block := range blockchain.Chain {
		fmt.Printf("{ \n index: %d, \n", block.index)
		fmt.Printf(" nonce: %d, \n", block.nonce)
		fmt.Printf(" timestamp: %d, \n", block.timestamp)
		fmt.Printf(" previousHash: %s, \n", block.PreviousHash)
		fmt.Printf(" currentHash: %s, \n } \n", block.CurrentHash)
	}
}

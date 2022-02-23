package Controllers

import (
	"blockchain/Models"
	"fmt"
)

func IsChainValid(blockchain *Models.Blockchain) bool {
	i := 0
	for i < blockchain.Length-1 {
		if !CheckBlock(blockchain.Chain[i]) {
			return false
		} else if fmt.Sprintf("%x", blockchain.Chain[i].CurrentHash) != fmt.Sprintf("%x", blockchain.Chain[i+1].PreviousHash) {
			return false
		}
		i += 1
	}
	return true
}

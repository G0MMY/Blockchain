package Controllers

import (
	"blockchain/Models"
	"fmt"
)

func IsChainValid(blockchain *Models.Blockchain) bool {
	i := 0
	for i < blockchain.Length-1 {
		if fmt.Sprintf("%x", Hash(blockchain.Chain[i])) != fmt.Sprintf("%x", blockchain.Chain[i+1].PreviousHash) {
			return false
		}
		i += 1
	}
	return true
}

package Controllers

import (
	"fmt"
)

func (blockchain *Models.Blockchain) IsChainValid() bool {
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
func (blockchain *Models.Blockchain) GetChain() []*Models.Block {
	return blockchain.Chain
}
func (blockchain *Models.Blockchain) GetLength() int {
	return blockchain.Length
}

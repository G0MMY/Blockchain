package Models

import (
	"crypto/sha256"
)

type Tree struct {
	RootNode *Node
}

type Node struct {
	Data      []byte
	LeftNode  *Node
	RightNode *Node
}

func CreateTree(transactions []*Transaction) *Tree {
	if len(transactions)%2 != 0 {
		transactions = append(transactions, transactions[len(transactions)-1])
	}
	var treeArray [][][]byte
	var row [][]byte

	for _, transaction := range transactions {
		row = append(row, transaction.Hash())
	}
	treeArray = append(treeArray, row)

	for true {
		var tempRow [][]byte
		if len(row) == 1 {
			break
		}

		i := 0
		for i <= len(row)-2 {
			tempRow = append(tempRow, appendHash(row[i], row[i+1]))
			i += 2
		}

		row = tempRow
		if len(row) != 1 && len(row)%2 != 0 {
			row = append(row, row[len(row)-1])
		}
		treeArray = append(treeArray, row)
	}

	return &Tree{}
}

func appendHash(first, second []byte) []byte {
	first = append(first, second...)

	hash := sha256.Sum256(first)

	return hash[:]
}

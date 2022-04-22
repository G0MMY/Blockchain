package Models

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
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
	if transactions == nil || len(transactions) == 0 {
		return nil
	}
	treeArray := createTreeArray(transactions)
	treeNode := linkNodes(len(treeArray)-1, 0, treeArray)

	return &Tree{treeNode}
}

func DecodeTree(byteTree []byte) *Tree {
	var tree Tree
	decoder := gob.NewDecoder(bytes.NewReader(byteTree))

	if err := decoder.Decode(&tree); err != nil {
		log.Panic(err)
	}

	return &tree
}

func (tree *Tree) EncoreTree() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	if err := encoder.Encode(tree); err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}

func (tree *Tree) CheckTree(transactions []*Transaction) bool {
	if len(transactions)%2 != 0 {
		transactions = append(transactions, transactions[len(transactions)-1])
	}

	var hash [][]byte
	for _, transaction := range transactions {
		hash = append(hash, transaction.Hash())
	}

	return browseTree(tree.RootNode, hash, 0)
}

func browseTree(node *Node, hash [][]byte, j int) bool {
	if j >= len(hash) {
		j -= 2
	}
	if node.LeftNode.Data == nil && node.RightNode.Data == nil {
		if bytes.Compare(node.Data, hash[j]) != 0 {
			return false
		}
		return true
	}
	if !browseTree(node.LeftNode, hash, j*2) {
		return false
	}
	if !browseTree(node.RightNode, hash, j*2+1) {
		return false
	}

	return true
}

func linkNodes(i, j int, treeArray [][][]byte) *Node {
	node := &Node{}
	if i >= 0 {
		if j >= len(treeArray[i]) {
			j -= 2
		}
		node.Data = treeArray[i][j]
		node.LeftNode = linkNodes(i-1, j*2, treeArray)
		node.RightNode = linkNodes(i-1, j*2+1, treeArray)
	}

	return node
}

func createTreeArray(transactions []*Transaction) [][][]byte {
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

	return treeArray
}

func appendHash(first, second []byte) []byte {
	first = append(first, second...)

	hash := sha256.Sum256(first)

	return hash[:]
}

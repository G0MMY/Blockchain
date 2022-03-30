package Models

type Node struct {
	Hash       []byte
	LeftChild  *Node
	RightChild *Node
	Leaf       *Transaction
}

type MerkleTree struct {
	root *Node
}

func (node Node) GetLeaf() Transaction {
	if node.LeftChild == nil && node.RightChild == nil {
		return *node.Leaf
	}

	return Transaction{}
}

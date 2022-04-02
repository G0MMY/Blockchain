package Models

type Tree struct {
	RootNode *Node
}

type Node struct {
	Data      []byte
	LeftNode  *Node
	RightNode *Node
}

func CreateTree(transactions []*Transaction) *Tree {
	return &Tree{}
}

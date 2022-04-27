package Network

type Handler struct {
	Node *FullNode
}

func New(node *FullNode) Handler {
	return Handler{node}
}

package Network

type HandlerNode struct {
	Node *FullNode
}

type HandlerMiner struct {
	Miner *Miner
}

func NewNode(node *FullNode) HandlerNode {
	return HandlerNode{node}
}

func NewMiner(miner *Miner) HandlerMiner {
	return HandlerMiner{miner}
}

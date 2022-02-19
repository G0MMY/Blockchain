package Models

type Blockchain struct {
	Chain  []*Block
	Length int
}

type IBlockchain interface {
	GetLength() int
	GetChain() []*Block
	IsChainValid() bool
}

package Models

type Transaction struct {
	ID      int
	BlockID int
	Block   Block
	Inputs  []Input
	Outputs []Output
}

type Input struct {
	ID            int
	TransactionId int
	OutputId      int
	Output        Output
	Signature     string
}

type Output struct {
	ID            int
	TransactionId int
	Amount        int
	PublicKey     []byte
}

type UnspentTransaction struct {
	ID      int
	Inputs  []Input
	Outputs []Output
	Fee     int
}

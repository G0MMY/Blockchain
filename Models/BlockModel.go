package Models

type Block struct {
	ID                    int
	Nonce                 int
	Timestamp             int64
	Transactions          []Transaction
	PreviousHash          []byte
	CurrentHash           []byte
	MaxNumberTransactions int
}

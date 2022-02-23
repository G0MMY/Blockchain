package Models

type Block struct {
	ID                    int `gorm:"autoIncrement"`
	Nonce                 int
	Timestamp             int64
	Transactions          []Transaction
	PreviousHash          []byte
	CurrentHash           []byte
	MaxNumberTransactions int
}

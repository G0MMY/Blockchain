package Models

type Transaction struct {
	ID        int `gorm:"autoIncrement"`
	BlockID   int
	Block     Block
	Inputs    []Input
	Outputs   []Output
	Timestamp int64
}

type Input struct {
	ID            int `gorm:"autoIncrement"`
	TransactionId int
	OutputId      int `gorm:"unique"`
	Output        Output
	Signature     string
}

type Output struct {
	ID            int `gorm:"autoIncrement"`
	TransactionId int
	Amount        int
	PublicKey     []byte
}

type MemPoolInput struct {
	ID                   int `gorm:"autoIncrement"`
	MemPoolTransactionId int
	OutputId             int `gorm:"unique"`
	Output               Output
	Signature            string
}

type MemPoolOutput struct {
	ID                   int `gorm:"autoIncrement"`
	MemPoolTransactionId int
	Amount               int
	PublicKey            []byte
}

type MemPoolTransaction struct {
	ID        int `gorm:"autoIncrement"`
	Inputs    []MemPoolInput
	Outputs   []MemPoolOutput
	Fee       int
	Timestamp int64
}

type CreateTransaction struct {
	Amount    int
	From      string
	To        string
	Signature string
	Fee       int
}

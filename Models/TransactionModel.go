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
	Signature     []byte
}

type Output struct {
	ID            int `gorm:"autoIncrement"`
	TransactionId int
	Amount        int
	PublicKey     string
}

type MemPoolInput struct {
	ID                   int `gorm:"autoIncrement"`
	MemPoolTransactionId int
	OutputId             int `gorm:"unique"`
	Output               Output
	Signature            []byte
}

type MemPoolOutput struct {
	ID                   int `gorm:"autoIncrement"`
	MemPoolTransactionId int
	Amount               int
	PublicKey            string
}

type MemPoolTransaction struct {
	ID        int `gorm:"autoIncrement"`
	Inputs    []MemPoolInput
	Outputs   []MemPoolOutput
	Fee       int
	Timestamp int64
}

//change that
type CreateTransaction struct {
	Amount     int
	To         string
	Fee        int
	PrivateKey string
}

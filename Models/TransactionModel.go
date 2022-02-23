package Models

type Transaction struct {
	ID      int `gorm:"autoIncrement"`
	BlockID int
	Block   Block
	Inputs  []Input
	Outputs []Output
}

type Input struct {
	ID            int `gorm:"autoIncrement"`
	TransactionId int
	Output        []Output
	Signature     string
}

type Output struct {
	ID            int `gorm:"autoIncrement"`
	TransactionId int
	InputId       int
	Amount        int
	PublicKey     []byte
}

type MemPoolInput struct {
	ID                   int `gorm:"autoIncrement"`
	MemPoolTransactionId int
	Output               []MemPoolOutput
	Signature            string
}

type MemPoolOutput struct {
	ID                   int `gorm:"autoIncrement"`
	MemPoolTransactionId int
	InputId              int
	Amount               int
	PublicKey            []byte
}

type MemPoolTransaction struct {
	ID      int `gorm:"autoIncrement"`
	Inputs  []MemPoolInput
	Outputs []MemPoolOutput
	Fee     int
}

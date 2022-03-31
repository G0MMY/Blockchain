package Models

type Transaction struct {
	ID        int `gorm:"autoIncrement"`
	BlockID   int
	Block     Block
	Inputs    []Input
	Outputs   []Output
	Timestamp int64
	Fee       int
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

func (transaction Transaction) TableName() string {
	if transaction.Fee != 0 && transaction.BlockID == 0 {
		return "memPool_Transaction"
	}

	return "transaction"
}

type CreateTransaction struct {
	Amount     int
	To         string
	Fee        int
	PrivateKey string
	Timestamp  int64
}

package Models

type Block struct {
	ID           int `gorm:"autoIncrement"`
	Nonce        int
	Timestamp    int64
	MerkleRoot   []byte
	PreviousHash []byte
	Difficulty   int
	Transactions []Transaction
}

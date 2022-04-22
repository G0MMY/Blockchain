package Models

type WalletResponse struct {
	PrivateKey string
	PublicKey  string
	Address    string
}

type CreateBlockRequest struct {
	Index        int
	Nonce        int
	Timestamp    int64
	MerkleRoot   string
	PreviousHash string
	Transactions []string
	MerkleTree   string
	PrivateKey   string
}

type CreateTransactionRequest struct {
	PrivateKey string
	To         string
	Amount     int
	Fee        int
	Timestamp  int64
}

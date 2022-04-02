package Models

import "time"

var (
	coinbaseAmount = 10
)

type Transaction struct {
	Inputs    []*Input
	Outputs   []*Output
	Timestamp int64
	Fee       int
}

type Input struct {
	OutputTransactionId []byte
	OutputIndex         int
	Signature           []byte
}

type Output struct {
	PublicKeyHash []byte
	Amount        int
}

//use pub key hash
func CreateCoinbase(address []byte) *Transaction {
	input := &Input{[]byte{}, -1, []byte{}}
	output := &Output{address, coinbaseAmount}

	return &Transaction{[]*Input{input}, []*Output{output}, time.Now().Unix(), 0}
}

package Models

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

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
	output    *Output
	Signature []byte
}

type Output struct {
	BlockId          []byte
	TransactionIndex int
	Index            int
	PublicKeyHash    []byte
	Amount           int
}

//use pub key hash
func CreateCoinbase(address []byte) *Transaction {
	input := &Input{&Output{}, []byte{}}
	output := &Output{address, coinbaseAmount}

	return &Transaction{[]*Input{input}, []*Output{output}, time.Now().Unix(), 0}
}

func (transaction *Transaction) IsCoinbase() bool {
	if transaction.Fee == 0 && len(transaction.Inputs) == 1 && len(transaction.Outputs) == 1 && transaction.Inputs[0].OutputIndex == -1 {
		return true
	}

	return false
}

func (transaction *Transaction) EncodeTransaction() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	if err := encoder.Encode(transaction); err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}

func DecodeTransaction(byteTransaction []byte) *Transaction {
	var transaction Transaction
	decoder := gob.NewDecoder(bytes.NewReader(byteTransaction))

	if err := decoder.Decode(&transaction); err != nil {
		log.Panic(err)
	}

	return &transaction
}

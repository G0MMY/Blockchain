package Models

import (
	"bytes"
	"crypto/sha256"
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

type UnspentOutput struct {
	outputs []*Output
}

//use pub key hash
func CreateCoinbase(address []byte) *Transaction {
	input := &Input{&Output{}, []byte{}}
	output := &Output{[]byte{}, -1, 0, address, coinbaseAmount}

	return &Transaction{[]*Input{input}, []*Output{output}, time.Now().Unix(), 0}
}

func CreateTransaction(to, from []byte, amount, amountRest, fee int, timestamp int64, unspentOutputs *UnspentOutput) *Transaction {
	var outputs []*Output
	inputs := unspentOutputs.CreateInputs()

	if amountRest == 0 {
		outputs = append(outputs, &Output{[]byte{}, -1, -1, to, amount})
		outputs = append(outputs, &Output{[]byte{}, -1, -1, from, amountRest})
	}

	return &Transaction{inputs, outputs, timestamp, fee}
}

func (transaction *Transaction) Hash() []byte {
	hash := sha256.Sum256(transaction.EncodeTransaction())

	return hash[:]
}

func (transaction *Transaction) IsCoinbase() bool {
	if transaction.Fee == 0 && len(transaction.Inputs) == 1 && len(transaction.Outputs) == 1 && transaction.Inputs[0].output == nil {
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

//add signature
func (unspentOutputs *UnspentOutput) CreateInputs() []*Input {
	var inputs []*Input

	for _, output := range unspentOutputs.outputs {
		inputs = append(inputs, &Input{output, []byte{}})
	}

	return inputs
}

func (unspentOutputs *UnspentOutput) GetOutputsForAmount(amount int) ([]*Output, int) {
	var outputs []*Output
	index := -1

	for i, output := range unspentOutputs.outputs {
		if amount > 0 {
			outputs = append(outputs, output)
			amount -= output.Amount
		} else if index == -1 {
			index = i
		}
	}
	rest := unspentOutputs.outputs[index:]
	unspentOutputs.outputs = outputs

	return rest, amount
}

func (unspentOutput *UnspentOutput) EncodeUnspentOutput() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	if err := encoder.Encode(unspentOutput); err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}

func DecodeUnspentOutput(byteOutput []byte) *UnspentOutput {
	var output UnspentOutput
	decoder := gob.NewDecoder(bytes.NewReader(byteOutput))

	if err := decoder.Decode(&output); err != nil {
		log.Panic(err)
	}

	return &output
}



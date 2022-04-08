package Models

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

var (
	coinbaseAmount = 50
)

type Transaction struct {
	Inputs    []*Input
	Outputs   []*Output
	Timestamp int64
	Fee       int
}

type Input struct {
	Output    *Output
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
	Outputs []*Output
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

	if amountRest < 0 {
		outputs = append(outputs, &Output{[]byte{}, -1, -1, to, amount})
		outputs = append(outputs, &Output{[]byte{}, -1, -1, from, amountRest * -1})
	} else if amountRest == 0 {
		outputs = append(outputs, &Output{[]byte{}, -1, -1, to, amount})
	}

	return &Transaction{inputs, outputs, timestamp, fee}
}

func insertTransaction(transactions []*Transaction, transaction *Transaction) []*Transaction {
	for i, trans := range transactions {
		if transaction.Fee > trans.Fee {
			transactions = append(transactions[:i+1], transactions[i:]...)
			transactions[i] = transaction
			break
		}
	}

	return transactions
}

//add validation
func FindBestMemPoolTransactions(transactions []*Transaction, numberTransactions int) []*Transaction {
	var memPoolTransactions []*Transaction
	if len(transactions) > 0 {
		memPoolTransactions = append(memPoolTransactions, transactions[0])
		i := 1

		for i < len(transactions) {
			if transactions[i].Timestamp < time.Now().Unix() {
				//if ValidateTransaction(transactions[i]) {
				if len(memPoolTransactions) < numberTransactions {
					for j, transaction := range memPoolTransactions {
						if transaction.Fee <= transactions[i].Fee {
							memPoolTransactions = append(memPoolTransactions[:j+1], memPoolTransactions[j:]...)
							memPoolTransactions[j] = transactions[i]
							break
						}
					}
				} else if memPoolTransactions[len(memPoolTransactions)-1].Fee < transactions[i].Fee {
					memPoolTransactions = insertTransaction(memPoolTransactions, transactions[i])[:numberTransactions]
				}
				//}
			}
			i++
		}
	}

	return memPoolTransactions
}

func (transaction *Transaction) Hash() []byte {
	hash := sha256.Sum256(transaction.EncodeTransaction())

	return hash[:]
}

func (transaction *Transaction) IsCoinbase() bool {
	if transaction.Fee == 0 && len(transaction.Inputs) == 1 && len(transaction.Outputs) == 1 && transaction.Inputs[0].Output.PublicKeyHash == nil {
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

	for _, output := range unspentOutputs.Outputs {
		inputs = append(inputs, &Input{output, []byte{}})
	}

	return inputs
}

func (unspentOutputs *UnspentOutput) GetOutputsForAmount(amount int) ([]*Output, int) {
	var outputs []*Output
	rest := unspentOutputs.Outputs
	index := -1

	for i, output := range unspentOutputs.Outputs {
		if amount > 0 {
			outputs = append(outputs, output)
			amount -= output.Amount
		} else if index == -1 {
			index = i
		}
	}
	unspentOutputs.Outputs = outputs

	if index == -1 {
		return nil, amount
	}

	return rest[index:], amount
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

func ValidateAddress(address []byte) []byte {
	if IsValidPublicKey(address) {
		return GetPublicKeyHash(address)
	} else if IsValidAddress(address) {
		return GetPublicKeyHashFromAddress(address)
	} else {
		log.Panic("Invalid address provided")
	}

	return []byte{}
}

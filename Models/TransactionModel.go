package Models

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"sort"
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
	PublicKey []byte
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

func CreateCoinbase(privateKey []byte) *Transaction {
	input := &Input{&Output{Amount: coinbaseAmount}, Sign(coinbaseAmount, privateKey), GetPublicKeyFromPrivateKey(privateKey)}
	output := &Output{[]byte{}, -1, 0, ValidateAddress(GetPublicKeyFromPrivateKey(privateKey)), coinbaseAmount}

	return &Transaction{[]*Input{input}, []*Output{output}, time.Now().Unix(), 0}
}

func CreateTransaction(to, from, privateKey []byte, amount, amountRest, fee int, timestamp int64, unspentOutputs *UnspentOutput) *Transaction {
	var outputs []*Output
	inputs := unspentOutputs.CreateInputs(privateKey)

	if amountRest < 0 {
		outputs = append(outputs, &Output{[]byte{}, -1, -1, to, amount})
		outputs = append(outputs, &Output{[]byte{}, -1, -1, from, amountRest * -1})
	} else if amountRest == 0 {
		outputs = append(outputs, &Output{[]byte{}, -1, -1, to, amount})
	}

	return &Transaction{inputs, outputs, timestamp, fee}
}

func FindBestMemPoolTransactions(transactions []*Transaction, numberTransactions int, privateKey []byte) ([]*Transaction, [][]byte) {
	if !IsValidPrivateKey(privateKey) {
		log.Panic("Invalid private key")
	}
	var memPoolTransactions []*Transaction
	var transactionsHash [][]byte

	if len(transactions) > 0 {
		sort.Slice(transactions, func(i, j int) bool {
			return transactions[i].Fee > transactions[j].Fee
		})

		i := 0
		for i < len(transactions) {
			if transactions[i].Timestamp <= time.Now().Unix() {
				if len(memPoolTransactions) < numberTransactions-1 {
					transactionsHash = append(transactionsHash, transactions[i].Hash())
					transactions[i].addFeeOutput(privateKey)
					transactions[i].ValidateTransaction()
					memPoolTransactions = append(memPoolTransactions, transactions[i])
				} else {
					break
				}
			}
			i += 1
		}
	}

	return memPoolTransactions, transactionsHash
}

func (transaction *Transaction) addFeeOutput(privateKey []byte) {
	if !IsValidPrivateKey(privateKey) {
		log.Panic("Invalid private key")
	}
	transaction.Outputs = append(transaction.Outputs, &Output{[]byte{}, -1, -1, ValidateAddress(GetPublicKeyFromPrivateKey(privateKey)), transaction.Fee})
}

func (transaction *Transaction) ValidateTransaction() {
	inputAmount := 0
	for _, input := range transaction.Inputs {
		inputAmount += input.Output.Amount
		if !ValidateSignature(input.Output.Amount, input.PublicKey, input.Signature) {
			log.Panic("Invalid signature")
		}
	}

	outputAmount := 0
	for _, output := range transaction.Outputs {
		outputAmount += output.Amount
	}

	if inputAmount != outputAmount {
		log.Panic("Not all money is there")
	}
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

func (unspentOutputs *UnspentOutput) CreateInputs(privateKey []byte) []*Input {
	var inputs []*Input

	for _, output := range unspentOutputs.Outputs {
		inputs = append(inputs, &Input{output, Sign(output.Amount, privateKey), GetPublicKeyFromPrivateKey(privateKey)})
	}

	return inputs
}

func (unspentOutputs *UnspentOutput) GetOutputsForAmount(amount int) ([]*Output, int) {
	var outputs []*Output

	if len(unspentOutputs.Outputs) == 0 {
		log.Panic("No unspent outputs to choose from")
	}

	sort.Slice(unspentOutputs.Outputs, func(i, j int) bool {
		return unspentOutputs.Outputs[i].Amount > unspentOutputs.Outputs[j].Amount
	})

	rest := unspentOutputs.Outputs

	if amount > unspentOutputs.Outputs[0].Amount {
		index := -1
		for i, output := range unspentOutputs.Outputs {
			if amount > 0 {
				outputs = append(outputs, output)
				amount -= output.Amount
			} else if index == -1 {
				index = i
				break
			}
		}
		if index == -1 && amount > 0 {
			return nil, amount
		}
		rest = rest[index:]
	} else {
		for i, output := range unspentOutputs.Outputs {
			if amount < output.Amount && i < len(unspentOutputs.Outputs)-1 {
				if amount > unspentOutputs.Outputs[i+1].Amount {
					outputs = append(outputs, output)
					amount -= output.Amount
					rest = append(rest[:i], rest[i+1:]...)
					break
				}
			} else {
				outputs = append(outputs, output)
				amount -= output.Amount
				rest = append(rest[:i], rest[i+1:]...)
				break
			}
		}
	}

	if len(rest) == len(unspentOutputs.Outputs) {
		return nil, amount
	}

	unspentOutputs.Outputs = outputs

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

func GenerateUnspentOutputKey(address []byte) []byte {
	return bytes.Join([][]byte{
		[]byte("UnspentOutput-"),
		address,
	}, []byte{})
}

func DecodeUnspentOutput(byteOutput []byte) *UnspentOutput {
	var output UnspentOutput
	decoder := gob.NewDecoder(bytes.NewReader(byteOutput))

	if err := decoder.Decode(&output); err != nil {
		log.Panic(err)
	}

	return &output
}

package Models

import (
	"bytes"
	"crypto/sha256"
	"github.com/ugorji/go/codec"
	"log"
	"sort"
	"time"
)

var (
	MemPoolPrefix       = []byte("MemPool-")
	UnspentOutputPrefix = []byte("UnspentOutput-")
	coinbaseAmount      = 50
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
	PublicKeyHash []byte
	Amount        int
}

type UnspentOutput struct {
	Outputs []*Output
}

//maybe check that not sure if its good
func CreateCoinbase(privateKey []byte) *Transaction {
	input := &Input{&Output{Amount: coinbaseAmount}, Sign(coinbaseAmount, privateKey), GetPublicKeyFromPrivateKey(privateKey)}
	output := &Output{ValidateAddress(GetPublicKeyFromPrivateKey(privateKey)), coinbaseAmount}

	return &Transaction{[]*Input{input}, []*Output{output}, time.Now().Unix(), 0}
}

func CreateTransaction(to, from, privateKey []byte, amount, amountRest, fee int, timestamp int64, unspentOutputs *UnspentOutput) *Transaction {
	var outputs []*Output
	inputs := unspentOutputs.CreateInputs(privateKey)

	if amountRest < 0 {
		outputs = append(outputs, &Output{to, amount})
		outputs = append(outputs, &Output{from, amountRest * -1})
	} else if amountRest == 0 {
		outputs = append(outputs, &Output{to, amount})
	}

	return &Transaction{inputs, outputs, timestamp, fee}
}

func FindBestMemPoolTransactions(transactions []*Transaction, numberTransactions int, privateKey []byte) []*Transaction {
	if !IsValidPrivateKey(privateKey) {
		log.Println("Invalid private key")
		return nil
	}
	var memPoolTransactions []*Transaction

	if len(transactions) > 0 {
		sort.Slice(transactions, func(i, j int) bool {
			return transactions[i].Fee > transactions[j].Fee
		})

		i := 0
		for i < len(transactions) {
			if transactions[i].Timestamp <= time.Now().Unix() {
				if len(memPoolTransactions) < numberTransactions-1 {
					if !transactions[i].addFeeOutput(privateKey) || !transactions[i].ValidateTransaction(false) {
						return nil
					}
					memPoolTransactions = append(memPoolTransactions, transactions[i])
				} else {
					break
				}
			}
			i += 1
		}
	}

	return memPoolTransactions
}

func HashUnspentOutputs(unspentOutputs map[string]*UnspentOutput) []byte {
	var byteUnspentOutputs [][]byte
	for _, unspentOutput := range unspentOutputs {
		byteUnspentOutputs = append(byteUnspentOutputs, unspentOutput.Hash())
	}

	hash := sha256.Sum256(bytes.Join(byteUnspentOutputs, []byte{}))

	return hash[:]
}

func HashTransactions(transactions []*Transaction) []byte {
	var byteTransactions [][]byte

	for _, transaction := range transactions {
		hashTransaction := transaction.Hash()
		if hashTransaction == nil {
			return nil
		}
		byteTransactions = append(byteTransactions, hashTransaction)
	}

	hash := sha256.Sum256(bytes.Join(byteTransactions, []byte{}))

	return hash[:]
}

func (transaction *Transaction) GetMemPoolHash(block *Block) []byte {
	temp := *transaction
	var tempOutputs []*Output
	for _, output := range temp.Outputs {
		if output.Amount == temp.Fee && bytes.Compare(output.PublicKeyHash, block.Transactions[len(block.Transactions)-1].Outputs[0].PublicKeyHash) == 0 {
		} else {
			tempOutput := *output
			tempOutputs = append(tempOutputs, &tempOutput)
		}
	}
	temp.Outputs = tempOutputs

	return temp.Hash()
}

func (transaction *Transaction) addFeeOutput(privateKey []byte) bool {
	if !IsValidPrivateKey(privateKey) {
		log.Println("Invalid private key")
		return false
	}
	transaction.Outputs = append(transaction.Outputs, &Output{ValidateAddress(GetPublicKeyFromPrivateKey(privateKey)), transaction.Fee})

	return true
}

func (transaction *Transaction) ValidateTransaction(isMemPool bool) bool {
	inputAmount := 0
	for _, input := range transaction.Inputs {
		inputAmount += input.Output.Amount
		if !ValidateSignature(input.Output.Amount, input.PublicKey, input.Signature) {
			log.Println("Invalid signature")

			return false
		}
	}

	outputAmount := 0
	for _, output := range transaction.Outputs {
		outputAmount += output.Amount
	}

	if !isMemPool && inputAmount != outputAmount {
		log.Println("Not all money is there")

		return false
	} else if isMemPool && inputAmount-outputAmount != transaction.Fee && !transaction.IsCoinbase() {
		log.Println("The transaction output's are broken")

		return false
	}

	return true
}

func (transaction *Transaction) Hash() []byte {
	byteTransaction := transaction.EncodeTransaction()
	if byteTransaction == nil {
		return nil
	}
	hash := sha256.Sum256(byteTransaction)

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
	encoder := codec.NewEncoder(&buffer, new(codec.JsonHandle))

	if err := encoder.Encode(transaction.CreateTransactionRequest()); err != nil {
		log.Println(err)
		return nil
	}

	return buffer.Bytes()
}

func DecodeTransaction(byteTransaction []byte) *Transaction {
	var transaction TransactionRequest
	decoder := codec.NewDecoder(bytes.NewReader(byteTransaction), new(codec.JsonHandle))

	if err := decoder.Decode(&transaction); err != nil {
		log.Println(err)
		return nil
	}

	return transaction.CreateTransaction()
}

func (unspentOutput *UnspentOutput) CreateInputs(privateKey []byte) []*Input {
	var inputs []*Input

	for _, output := range unspentOutput.Outputs {
		inputs = append(inputs, &Input{output, Sign(output.Amount, privateKey), GetPublicKeyFromPrivateKey(privateKey)})
	}

	return inputs
}

func (unspentOutput *UnspentOutput) GetOutputsForAmount(amount int) ([]*Output, int) {
	var outputs []*Output

	if len(unspentOutput.Outputs) == 0 {
		log.Println("No unspent outputs to choose from")
	}

	sort.Slice(unspentOutput.Outputs, func(i, j int) bool {
		return unspentOutput.Outputs[i].Amount > unspentOutput.Outputs[j].Amount
	})

	rest := unspentOutput.Outputs

	if amount > unspentOutput.Outputs[0].Amount {
		index := -1
		for i, output := range unspentOutput.Outputs {
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
		for i, output := range unspentOutput.Outputs {
			if amount < output.Amount && i < len(unspentOutput.Outputs)-1 {
				if amount > unspentOutput.Outputs[i+1].Amount {
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

	if len(rest) == len(unspentOutput.Outputs) {
		return nil, amount
	}

	unspentOutput.Outputs = outputs

	return rest, amount
}

func (unspentOutput *UnspentOutput) Hash() []byte {
	hash := sha256.Sum256(unspentOutput.EncodeUnspentOutput())

	return hash[:]
}

func (unspentOutput *UnspentOutput) EncodeUnspentOutput() []byte {
	var buffer bytes.Buffer
	encoder := codec.NewEncoder(&buffer, new(codec.JsonHandle))

	if err := encoder.Encode(unspentOutput.CreateUnspentOutputRequest()); err != nil {
		log.Println(err)
	}

	return buffer.Bytes()
}

func DecodeUnspentOutput(byteOutput []byte) *UnspentOutput {
	var unspentOutput UnspentOutputsRequest
	decoder := codec.NewDecoder(bytes.NewReader(byteOutput), new(codec.JsonHandle))

	if err := decoder.Decode(&unspentOutput); err != nil {
		log.Println(err)
	}

	return unspentOutput.CreateUnspentOutput()
}

func GenerateUnspentOutputKey(address []byte) []byte {
	return bytes.Join([][]byte{
		UnspentOutputPrefix,
		address,
	}, []byte{})
}

func GenerateMemPoolTransactionKey(hash []byte) []byte {
	return bytes.Join([][]byte{
		MemPoolPrefix,
		hash,
	}, []byte{})
}

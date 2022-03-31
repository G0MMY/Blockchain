package Controllers

import (
	"blockchain/Models"
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"
)

func CreateMemPoolTransaction(inputs []Models.Input, outputs []Models.Output, fee int, timestamp int64) Models.Transaction {
	if fee <= 0 {
		panic("The fee must be greater then 0")
		return Models.Transaction{}
	}

	return Models.Transaction{Inputs: inputs, Outputs: outputs, Timestamp: timestamp, Fee: fee}
}

func CreateMemPoolInputs(outputs []Models.Output) []Models.Input {
	var inputs []Models.Input
	for _, output := range outputs {
		inputs = append(inputs, Models.Input{OutputId: output.ID, Output: output})
	}
	return inputs
}

func CreateMemPoolOutputs(amount int, to string, inputs []Models.Input) []Models.Output {
	totalAmount := 0
	for _, input := range inputs {
		totalAmount += input.Output.Amount
	}

	outputs := append([]Models.Output{}, Models.Output{Amount: amount, PublicKey: to})
	if totalAmount-amount > 0 {
		outputs = append(outputs, Models.Output{Amount: totalAmount - amount, PublicKey: inputs[0].Output.PublicKey})
	}

	return outputs
}

func LinkInputs(inputs []Models.Input, outputs []Models.Output) []Models.Input {
	result := inputs

	for i, input := range inputs {
		for _, output := range outputs {
			if output.ID == input.OutputId {
				result[i].Output = output
				break
			}
		}
	}

	return result
}

func LinkTransactions(transactions []Models.Transaction, inputs []Models.Input, outputs []Models.Output) []Models.Transaction {
	i := 0

	for i < len(transactions) {
		for _, input := range inputs {
			if input.TransactionId == transactions[i].ID {
				transactions[i].Inputs = append(transactions[i].Inputs, input)
			}
		}
		for _, output := range outputs {
			if output.TransactionId == transactions[i].ID {
				transactions[i].Outputs = append(transactions[i].Outputs, output)
			}
		}
		i += 1
	}

	return transactions
}

func BuildTransaction(outputs []Models.Output, body Models.CreateTransaction, privateKey []byte) Models.Transaction {
	memPoolInput := CreateMemPoolInputs(outputs)
	memPoolOutput := CreateMemPoolOutputs(body.Amount, body.To, memPoolInput)
	memPoolTransaction := CreateMemPoolTransaction(memPoolInput, memPoolOutput, body.Fee, body.Timestamp)
	return SignTransaction(privateKey, memPoolTransaction)
}

func GetMemPoolTransactions(transactions []Models.Transaction) []Models.Transaction {
	var result []Models.Transaction

	for _, transaction := range transactions {
		if transaction.TableName() == "memPool_Transaction" {
			result = append(result, transaction)
		}
	}

	return result
}

func GetOutputs(outputs []Models.Output, amount int) []Models.Output {
	var result []Models.Output
	for _, output := range outputs {
		result = append(result, output)
		amount -= output.Amount
		if amount <= 0 {
			break
		}
	}
	if amount > 0 {
		return nil
	}
	return result
}

func FindBestMemPoolTransactions(transactions []Models.Transaction, numberTransactions int) []Models.Transaction {
	var memPoolTransactions []Models.Transaction
	i := 0

	for i < len(transactions) {
		if transactions[i].Timestamp < time.Now().Unix() && transactions[i].TableName() == "memPool_Transaction" {
			if ValidateTransaction(transactions[i]) {
				if len(memPoolTransactions) < numberTransactions {
					if len(memPoolTransactions) == 0 {
						memPoolTransactions = append(memPoolTransactions, transactions[i])
					} else {
						for j, transaction := range memPoolTransactions {
							if transaction.Fee <= transactions[i].Fee {
								memPoolTransactions = append(memPoolTransactions[:j+1], memPoolTransactions[j:]...)
								memPoolTransactions[j] = transactions[i]
								break
							}
						}
					}
				} else if memPoolTransactions[len(memPoolTransactions)-1].Fee < transactions[i].Fee {
					memPoolTransactions = insertTransaction(memPoolTransactions, transactions[i])[:numberTransactions]
				}
			}
		}
		i++
	}

	return memPoolTransactions
}

func insertTransaction(transactions []Models.Transaction, transaction Models.Transaction) []Models.Transaction {
	for i, trans := range transactions {
		if transaction.Fee > trans.Fee {
			transactions = append(transactions[:i+1], transactions[i:]...)
			transactions[i] = transaction
			break
		}
	}

	return transactions
}

func GetMemPoolTransactionsIds(memPoolTransactions []Models.Transaction) []int {
	var ids []int
	for _, transaction := range memPoolTransactions {
		ids = append(ids, transaction.ID)
	}
	return ids
}

func TransactionsToByte(transactions []Models.Transaction) []byte {
	var byteArray [][]byte
	for _, transaction := range transactions {
		byteArray = append(byteArray, []byte(fmt.Sprintf("%x", transaction)))
	}
	return bytes.Join(byteArray, []byte{})
}

func hashTransaction(transaction Models.Transaction) []byte {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%x", transaction)))

	return hash[:]
}

func buildMerkleTree(transactions []Models.Transaction) Models.MerkleTree {
	
}

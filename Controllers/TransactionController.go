package Controllers

import (
	"blockchain/Models"
	"bytes"
	"fmt"
	"time"
)

func CreateMemPoolTransaction(inputs []Models.MemPoolInput, outputs []Models.MemPoolOutput, fee int) Models.MemPoolTransaction {
	return Models.MemPoolTransaction{Inputs: inputs, Outputs: outputs, Timestamp: time.Now().Unix(), Fee: fee}
}

//check signature
func CreateMemPoolInputs(signature string, outputs []Models.Output) []Models.MemPoolInput {
	var inputs []Models.MemPoolInput
	for _, output := range outputs {
		inputs = append(inputs, Models.MemPoolInput{OutputId: output.ID, Output: output, Signature: signature})
	}
	return inputs
}

func CreateMemPoolOutputs(amount int, to []byte, inputs []Models.MemPoolInput) []Models.MemPoolOutput {
	totalAmount := 0
	for _, input := range inputs {
		totalAmount += input.Output.Amount
	}
	var outputs []Models.MemPoolOutput
	outputs = append(outputs, Models.MemPoolOutput{Amount: amount, PublicKey: to})
	if totalAmount-amount > 0 {
		outputs = append(outputs, Models.MemPoolOutput{Amount: totalAmount - amount, PublicKey: inputs[0].Output.PublicKey})
	}
	return outputs
}

func CreateTransactions(memPoolTransactions []Models.MemPoolTransaction) []Models.Transaction {
	var transactions []Models.Transaction

	for _, memPoolTransaction := range memPoolTransactions {
		transactions = append(transactions, createTransaction(memPoolTransaction))
	}
	return transactions
}

func createTransaction(memPooltransaction Models.MemPoolTransaction) Models.Transaction {
	var inputs []Models.Input
	var outputs []Models.Output
	for _, input := range memPooltransaction.Inputs {
		inputs = append(inputs, memPoolInputToInput(input))
	}
	for _, output := range memPooltransaction.Outputs {
		outputs = append(outputs, memPoolOutputToInput(output))
	}
	return Models.Transaction{Inputs: inputs, Outputs: outputs, Timestamp: memPooltransaction.Timestamp}
}

func memPoolInputToInput(memPoolInput Models.MemPoolInput) Models.Input {
	return Models.Input{OutputId: memPoolInput.OutputId, Output: memPoolInput.Output, Signature: memPoolInput.Signature}
}

func memPoolOutputToInput(memPoolOutput Models.MemPoolOutput) Models.Output {
	return Models.Output{Amount: memPoolOutput.Amount, PublicKey: memPoolOutput.PublicKey}
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

//need to check timestamp
func FindBestMemPoolTransactions(transactions []Models.MemPoolTransaction, numberTransactions int) []Models.MemPoolTransaction {
	if len(transactions) < numberTransactions {
		return transactions[:len(transactions)]
	} else {
		result := transactions[:numberTransactions]
		if len(transactions) > numberTransactions {
			i := numberTransactions
			for i < len(transactions) {
				j := 0
				for j < numberTransactions {
					if transactions[i].Fee > result[j].Fee {
						result[j] = transactions[i]
						break
					}
					j += 1
				}
				i += 1
			}
		}
		return result
	}
}

func TransactionsToByte(transactions []Models.Transaction) []byte {
	var byteArray [][]byte
	for _, transaction := range transactions {
		byteArray = append(byteArray, []byte(fmt.Sprintf("%x", transaction)))
	}
	return bytes.Join(byteArray, []byte{})
}

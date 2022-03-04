package Handlers

import (
	"blockchain/Controllers"
	"blockchain/Models"
	"encoding/json"
	"fmt"
	"net/http"
)

func (h Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var body Models.CreateTransaction
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	defer r.Body.Close()

	outputs := Controllers.GetOutputs(h.getPublicKeyOutputs([]byte(body.From)), body.Amount)
	if outputs == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("you do not have enough money my friend")
		return
	}
	memPoolInput := Controllers.CreateMemPoolInputs(body.Signature, outputs)
	memPoolOutput := Controllers.CreateMemPoolOutputs(body.Amount, []byte(body.To), memPoolInput)
	memPoolTransaction := Controllers.CreateMemPoolTransaction(memPoolInput, memPoolOutput, body.Fee)

	if result := h.DB.Create(&memPoolTransaction); result.Error != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(result.Error)
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(memPoolTransaction)
	}
}

func (h Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	transactions := h.getTransactions()

	if transactions != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(transactions)
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("there are no transactions")
	}
}

func (h Handler) getTransactions() []Models.Transaction {
	var transactions []Models.Transaction
	outputs := h.getOutputs()
	inputs := h.getInputs()

	if result := h.DB.Find(&transactions); result.Error != nil {
		fmt.Println(result.Error)
	}

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

func (h Handler) getOutputs() []Models.Output {
	var outputs []Models.Output

	if result := h.DB.Find(&outputs); result.Error != nil {
		fmt.Println(result.Error)
	}
	return outputs
}

func (h Handler) getInputs() []Models.Input {
	var inputs []Models.Input

	if result := h.DB.Find(&inputs); result.Error != nil {
		fmt.Println(result.Error)
	}
	return inputs
}

func (h Handler) getPublicKeyOutputs(publicKey []byte) []Models.Output {
	var outputs []Models.Output

	result := h.DB.Where("public_key = ? "+
		"and outputs.id not in (select output_id from inputs) "+
		"and outputs.id not in (select output_id from mem_pool_inputs)", fmt.Sprintf("%s", publicKey)).Find(&outputs)

	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return outputs
}

func (h Handler) getMemPoolInputs() []Models.MemPoolInput {
	var inputs []Models.MemPoolInput

	if result := h.DB.Find(&inputs); result.Error != nil {
		fmt.Println(result.Error)
	}
	return inputs
}

func (h Handler) getMemPoolOutputs() []Models.MemPoolOutput {
	var outputs []Models.MemPoolOutput

	if result := h.DB.Find(&outputs); result.Error != nil {
		fmt.Println(result.Error)
	}
	return outputs
}

func (h Handler) GetMemPoolTransactions() []Models.MemPoolTransaction {
	var transactions []Models.MemPoolTransaction
	outputs := h.getMemPoolOutputs()
	inputs := h.getMemPoolInputs()

	if result := h.DB.Find(&transactions); result.Error != nil {
		fmt.Println(result.Error)
	}

	i := 0
	for i < len(transactions) {
		for _, input := range inputs {
			if input.MemPoolTransactionId == transactions[i].ID {
				transactions[i].Inputs = append(transactions[i].Inputs, input)
			}
		}
		for _, output := range outputs {
			if output.MemPoolTransactionId == transactions[i].ID {
				transactions[i].Outputs = append(transactions[i].Outputs, output)
			}
		}
		i += 1
	}
	return transactions
}

func (h Handler) DeleteMemPoolTransactions(ids []int) bool {
	var output Models.MemPoolOutput
	var input Models.MemPoolInput
	var transaction Models.MemPoolTransaction

	if result := h.DB.Where("mem_pool_transaction_id in ?", ids).Delete(&output); result.Error != nil {
		fmt.Println(result.Error)
		return false
	}
	if result := h.DB.Where("mem_pool_transaction_id in ?", ids).Delete(&input); result.Error != nil {
		fmt.Println(result.Error)
		return false
	}
	if result := h.DB.Where("id in ?", ids).Delete(&transaction); result.Error != nil {
		fmt.Println(result.Error)
		return false
	}
	return true
}

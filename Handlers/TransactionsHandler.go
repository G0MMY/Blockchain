package Handlers

import (
	"blockchain/Controllers"
	"blockchain/Models"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

//gotta check in memPool to check future outputs to do more transactions
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

	priv, pub := Controllers.GetDecodedKey(Controllers.StringKeyToByte(body.PrivateKey))
	outputs := Controllers.GetOutputs(h.getPublicKeyOutputs(pub), body.Amount)
	if outputs == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("you do not have enough money my friend")
		return
	}

	memPoolTransaction := Controllers.BuildTransaction(outputs, body, priv)
	if Controllers.ValidateTransaction(memPoolTransaction) {
		if result := h.DB.Create(&memPoolTransaction); result.Error != nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(result.Error)
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(memPoolTransaction)
		}
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Transaction not properly signed")
	}
}

func (h Handler) ValidateTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transaction := h.getTransaction(vars["transactionId"])

	if transaction.Inputs == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("The id dosen't belong to any transaction")
	} else {
		ok := Controllers.ValidateTransaction(transaction)

		if ok {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode("The transaction is validated")
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Transaction not properly signed")
		}
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

func (h Handler) getTransaction(id string) Models.Transaction {
	var transaction Models.Transaction
	outputs := h.getOutputsId(id)
	inputs := h.getInputsId(id)
	inputs = Controllers.LinkInputs(inputs, h.getOutputs())

	if result := h.DB.Find(&transaction, id); result.Error != nil {
		fmt.Println(result.Error)
	}

	return Controllers.LinkTransactions([]Models.Transaction{transaction}, inputs, outputs)[0]
}

func (h Handler) getTransactions() []Models.Transaction {
	var transactions []Models.Transaction
	outputs := h.getOutputs()
	inputs := h.getInputs()

	if result := h.DB.Find(&transactions); result.Error != nil {
		fmt.Println(result.Error)
	}

	return Controllers.LinkTransactions(transactions, inputs, outputs)
}

func (h Handler) getOutputs() []Models.Output {
	var outputs []Models.Output

	if result := h.DB.Find(&outputs); result.Error != nil {
		fmt.Println(result.Error)
	}
	return outputs
}

func (h Handler) getOutputsId(id string) []Models.Output {
	var outputs []Models.Output

	if result := h.DB.Where("transaction_id = ?", id).Find(&outputs); result.Error != nil {
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

func (h Handler) getInputsId(id string) []Models.Input {
	var inputs []Models.Input

	if result := h.DB.Where("transaction_id = ?", id).Find(&inputs); result.Error != nil {
		fmt.Println(result.Error)
	}
	return inputs
}

func (h Handler) getPublicKeyOutputs(publicKey []byte) []Models.Output {
	var outputs []Models.Output

	resultOutputs := h.DB.Where("public_key = ? "+
		"and outputs.id not in (select output_id from inputs)", hex.EncodeToString(publicKey)).Find(&outputs)

	if resultOutputs.Error != nil {
		fmt.Println(resultOutputs.Error)
	}
	return outputs
}

func (h Handler) getMemPoolInputs() []Models.Input {
	var inputs []Models.Input

	if result := h.DB.Find(&inputs); result.Error != nil {
		fmt.Println(result.Error)
	}

	for i, input := range inputs {
		var output Models.Output

		if result := h.DB.Where("id = ?", input.OutputId).Find(&output); result.Error != nil {
			fmt.Println(result.Error)
			break
		}

		inputs[i].Output = output
	}

	return inputs
}

func (h Handler) getMemPoolOutputs() []Models.Output {
	var outputs []Models.Output

	if result := h.DB.Find(&outputs); result.Error != nil {
		fmt.Println(result.Error)
	}
	return outputs
}

func (h Handler) GetMemPoolTransactions() []Models.Transaction {
	var transactions []Models.Transaction
	outputs := h.getMemPoolOutputs()
	inputs := h.getMemPoolInputs()

	if result := h.DB.Find(&transactions); result.Error != nil {
		fmt.Println(result.Error)
	}

	return Controllers.GetMemPoolTransactions(Controllers.LinkTransactions(transactions, inputs, outputs))
}

//func (h Handler) DeleteMemPoolTransactions(ids []int) bool {
//	var output Models.MemPoolOutput
//	var input Models.MemPoolInput
//	var transaction Models.MemPoolTransaction
//
//	if result := h.DB.Where("mem_pool_transaction_id in ?", ids).Delete(&output); result.Error != nil {
//		fmt.Println(result.Error)
//		return false
//	}
//	if result := h.DB.Where("mem_pool_transaction_id in ?", ids).Delete(&input); result.Error != nil {
//		fmt.Println(result.Error)
//		return false
//	}
//	if result := h.DB.Where("id in ?", ids).Delete(&transaction); result.Error != nil {
//		fmt.Println(result.Error)
//		return false
//	}
//	return true
//}

func (h Handler) LinkTransactions(block Models.Block, ids []int) bool {
	if result := h.DB.Table("transaction").Where("id in ?", ids).Updates(map[string]interface{}{"block_id": block.ID, "fee": 0}); result.Error != nil {
		fmt.Println(result.Error)
		return false
	}

	return true
}

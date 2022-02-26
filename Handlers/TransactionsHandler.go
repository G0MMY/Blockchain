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

	outputs := Controllers.GetOutputs(h.getOutputs([]byte(body.From)), body.Amount)
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

func (h Handler) getOutputs(publickKey []byte) []Models.Output {
	var outputs []Models.Output

	if result := h.DB.Where("public_key = ?", fmt.Sprintf("%s", publickKey)).Find(&outputs); result.Error != nil {
		fmt.Println("error with outputs")
	}
	return outputs
}

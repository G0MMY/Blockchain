package Network

import (
	"blockchain/Models"
	"encoding/json"
	"fmt"
	"net/http"
)

func (handler *HandlerNode) GetMemPoolTransactionsHash(w http.ResponseWriter, r *http.Request) {
	if handler.Node.Blockchain.DB == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("The blockchain's DB is not initialized")
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fmt.Sprintf("%x", handler.Node.Blockchain.GetMemPoolTransactionsHash()))
	}
}

func (handler *HandlerNode) GetMemPoolTransactions(w http.ResponseWriter, r *http.Request) {
	if handler.Node.Blockchain.DB == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("The blockchain's DB is not initialized")
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(handler.Node.Blockchain.GetMemPoolTransactions())
	}
}

func (handler *HandlerNode) AddTransaction(w http.ResponseWriter, r *http.Request) {
	var body Models.CreateTransactionRequest
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	defer r.Body.Close()

	AddTransaction <- body
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Added transaction")
}

func (handler *HandlerNode) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var body Models.TransactionRequest
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	defer r.Body.Close()

	CreateTransaction <- body
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Created transaction")
}

func (handler *HandlerNode) GetAllUnspentOutputs(w http.ResponseWriter, r *http.Request) {
	if handler.Node.Blockchain.DB == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("The blockchain's DB is not initialized")
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(handler.Node.Blockchain.GetAllUnspentOutputs())
	}
}

func (handler *HandlerNode) GetAllUnspentOutputsHash(w http.ResponseWriter, r *http.Request) {
	if handler.Node.Blockchain.DB == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("The blockchain's DB is not initialized")
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fmt.Sprintf("%x", handler.Node.Blockchain.GetAllUnspentOutputsHash()))
	}
}

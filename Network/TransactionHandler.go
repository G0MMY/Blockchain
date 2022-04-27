package Network

import (
	"blockchain/Models"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

func (handler *Handler) GetMemPoolTransactionsHash(w http.ResponseWriter, r *http.Request) {
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

func (handler *Handler) GetMemPoolTransactions(w http.ResponseWriter, r *http.Request) {
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

func (handler *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var body Models.CreateTransactionRequest
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	defer r.Body.Close()

	priv, err := hex.DecodeString(body.PrivateKey)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error while trying to decode the private key")
		return
	}

	to, er := hex.DecodeString(body.To)
	if er != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error while trying to decode the receiver")
		return
	}

	transaction := handler.Node.Blockchain.CreateTransaction(priv, to, body.Amount, body.Fee, body.Timestamp)

	if transaction == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while creating the transaction")
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(transaction)
	}
}

func (handler *Handler) GetAllUnspentOutputs(w http.ResponseWriter, r *http.Request) {
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

func (handler *Handler) GetAllUnspentOutputsHash(w http.ResponseWriter, r *http.Request) {
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

package Handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (handler *Handler) GetMemPoolTransactionsHash(w http.ResponseWriter, r *http.Request) {
	if handler.Blockchain.DB == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("The blockchain's DB is not initialized")
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fmt.Sprintf("%x", handler.Blockchain.HashMemPoolTransactions()))
	}
}

func (handler *Handler) GetMemPoolTransactions(w http.ResponseWriter, r *http.Request) {
	if handler.Blockchain.DB == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("The blockchain's DB is not initialized")
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(handler.Blockchain.GetMemPoolTransactions())
	}
}

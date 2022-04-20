package Handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (handler *Handler) GetLastHash(w http.ResponseWriter, r *http.Request) {
	if handler.Blockchain.LastHash == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("LastBlock not initialized")
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fmt.Sprintf("%x", handler.Blockchain.LastHash))
	}
}

func (handler *Handler) GetChain(w http.ResponseWriter, r *http.Request) {
	if handler.Blockchain.DB == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("The DB is nil")
	} else if handler.Blockchain.LastHash == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("LastBlock not initialized")
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(handler.Blockchain.GetBlockchain())
	}
}

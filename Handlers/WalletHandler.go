package Handlers

import (
	"blockchain/Controllers"
	"encoding/json"
	"net/http"
)

func (h Handler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	wallet := Controllers.CreateWallet()

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(wallet)
}

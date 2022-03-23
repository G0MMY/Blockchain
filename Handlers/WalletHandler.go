package Handlers

import (
	"blockchain/Controllers"
	"encoding/json"
	"net/http"
)

func (h Handler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	wallet := Controllers.CreateWallet()

	if result := h.DB.Create(&wallet); result.Error != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(result.Error)
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(wallet)
	}
}

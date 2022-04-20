package Handlers

import (
	"blockchain/Models"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func (handler *Handler) GetPublicKeyBalance(w http.ResponseWriter, r *http.Request) {
	if handler.Blockchain.DB == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("The blockchain's DB is not initialized")
	}

	vars := mux.Vars(r)
	address, err := hex.DecodeString(vars["address"])

	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error while trying to decode the address")
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(handler.Blockchain.GetBalance(address))
}

func (handler Handler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	wallet := Models.CreateWallet()

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Models.WalletResponse{fmt.Sprintf("%x", wallet.PrivateKey), fmt.Sprintf("%x", wallet.PublicKey), fmt.Sprintf("%x", Models.GetAddress(wallet.PublicKey))})
}

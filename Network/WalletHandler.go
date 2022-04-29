package Network

import (
	"blockchain/Models"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func (handler *HandlerNode) GetPublicKeyBalance(w http.ResponseWriter, r *http.Request) {
	if handler.Node.Blockchain.DB == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("The blockchain's DB is not initialized")
		return
	}

	vars := mux.Vars(r)
	address, err := hex.DecodeString(vars["address"])

	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error while trying to decode the address")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(handler.Node.Blockchain.GetBalance(address))
}

func (handler *HandlerNode) CreateWallet(w http.ResponseWriter, r *http.Request) {
	wallet := Models.CreateWallet()

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Models.WalletResponse{fmt.Sprintf("%x", wallet.PrivateKey), fmt.Sprintf("%x", wallet.PublicKey), fmt.Sprintf("%x", Models.GetAddress(wallet.PublicKey))})
}

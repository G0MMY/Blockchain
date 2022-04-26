package Handlers

import (
	"blockchain/Models"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func (handler *Handler) GetBlockMerkleRoot(w http.ResponseWriter, r *http.Request) {
	if handler.Blockchain.DB == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("The blockchain's DB is not initialized")
		return
	}

	vars := mux.Vars(r)
	blockHash, err := hex.DecodeString(vars["blockHash"])

	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error while trying to decode the block hash")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fmt.Sprintf("%x", handler.Blockchain.GetMerkleRoot(blockHash)))
}

func (handler *Handler) GetBlockTransactions(w http.ResponseWriter, r *http.Request) {
	if handler.Blockchain.DB == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("The blockchain's DB is not initialized")
		return
	}

	vars := mux.Vars(r)
	blockHash, err := hex.DecodeString(vars["blockHash"])

	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error while trying to decode the block hash")
	} else {
		block := handler.Blockchain.GetBlock(blockHash)
		if block == nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Invalid block hash provided")
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(block.Transactions)
		}
	}
}

func (handler *Handler) GetBlock(w http.ResponseWriter, r *http.Request) {
	if handler.Blockchain.DB == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("The blockchain's DB is not initialized")
		return
	}

	vars := mux.Vars(r)
	blockHash, err := hex.DecodeString(vars["blockHash"])

	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error while trying to decode the block hash")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(handler.Blockchain.GetBlock(blockHash).CreateBlockRequest())
}

func (handler *Handler) CreateBlock(w http.ResponseWriter, r *http.Request) {
	var body Models.CreateBlockResponse
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

	block, errorMessage := handler.Blockchain.CreateBlock(priv, Models.CreateBlockToBlock(&body))

	if block == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMessage)
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(block)
	}
}

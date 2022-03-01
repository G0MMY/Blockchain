package Handlers

import (
	"blockchain/Controllers"
	"blockchain/Models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (h Handler) AddGenesisBlock(w http.ResponseWriter, r *http.Request) {
	if h.GetLength() > 0 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Can't add Genesis Block on top of existing blocks")
	} else {
		outputs := append([]Models.Output{}, Models.Output{Amount: 10000, PublicKey: []byte("maxim")})
		transactions := append([]Models.Transaction{}, Models.Transaction{Outputs: outputs, Timestamp: time.Now().Unix()})
		block := Controllers.CreateBlock([]byte{0}, transactions)

		if result := h.DB.Create(&block); result.Error != nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(result.Error)
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(block)
		}
	}
}

func (h Handler) AddBlock(w http.ResponseWriter, r *http.Request) {
	trans, ids := h.GetMemPoolTransactions()
	memPoolTransactions := Controllers.FindBestMemPoolTransactions(trans, 5)
	block := Controllers.CreateBlock(h.getPreviousHash(), Controllers.CreateTransactions(memPoolTransactions))

	if block.PreviousHash != nil {
		if result := h.DB.Create(&block); result.Error != nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(result.Error)
		} else {
			if h.DeleteMemPoolTransactions(ids) {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(block)
			} else {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("There was an error with the transactions")
			}
		}
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Can't add block on top of nothing")
	}
}

func (h Handler) getPreviousHash() []byte {
	return Controllers.Hash(h.GetLastBlock())
}

func (h Handler) GetLastBlock() *Models.Block {
	var block Models.Block

	if result := h.DB.Last(&block); result.Error != nil {
		fmt.Println(result.Error)
	}

	return &block
}

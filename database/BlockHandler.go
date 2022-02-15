package database

import (
	"blockchain/components"
	"encoding/json"
	"fmt"
	"net/http"
)

func (h Handler) AddGenesisBlock(w http.ResponseWriter, r *http.Request) {
	if h.GetLength() > 0 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Can't add Genesis Block on top of existing blocks")
	} else {
		block := components.CreateBlock(1, []byte{0})

		if result := h.DB.Create(&block); result.Error != nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(result.Error)
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode("Created")
		}
	}
}

func (h Handler) AddBlock(w http.ResponseWriter, r *http.Request) {
	block := components.CreateBlock(h.GetLength()+1, h.getPreviousHash())

	if block.PreviousHash != nil {
		if result := h.DB.Create(&block); result.Error != nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(result.Error)
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode("Created")
		}
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Can't add block on top of nothing")
	}
}

func (h Handler) CheckLastBlock(w http.ResponseWriter, r *http.Request) {
	block := h.GetLastBlock()

	if block.Id != 0 {
		if block.CheckBlock() {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode("Block is good")
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode("Block isn't good")
		}
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Block not found")
	}
}

func (h Handler) getPreviousHash() []byte {
	return h.GetLastBlock().CurrentHash
}

func (h Handler) GetLastBlock() *components.BlockType {
	var block components.BlockType

	if result := h.DB.Last(&block); result.Error != nil {
		fmt.Println(result.Error)
	}

	return &block
}

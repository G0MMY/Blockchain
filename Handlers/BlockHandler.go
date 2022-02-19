package Handlers

import (
	"blockchain/Controllers"
	"blockchain/Models"
	"encoding/json"
	"fmt"
	"net/http"
)

func (h Handler) AddGenesisBlock(w http.ResponseWriter, r *http.Request) {
	if h.GetLength() > 0 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Can't add Genesis IBlock on top of existing blocks")
	} else {
		block := Controllers.CreateBlock(1, []byte{0})

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
	block := Controllers.CreateBlock(h.GetLength()+1, h.getPreviousHash())

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

	if block.ID != 0 {
		if block.CheckBlock() {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode("IBlock is good")
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode("IBlock isn't good")
		}
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("IBlock not found")
	}
}

func (h Handler) getPreviousHash() []byte {
	return h.GetLastBlock().CurrentHash
}

func (h Handler) GetLastBlock() *Models.Block {
	var block Models.Block

	if result := h.DB.Last(&block); result.Error != nil {
		fmt.Println(result.Error)
	}

	return &block
}

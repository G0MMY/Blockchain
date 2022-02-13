package database

import (
	"blockchain/components"
	"encoding/json"
	"fmt"
	"net/http"
)

func (h handler) GetBlocks(w http.ResponseWriter, r *http.Request) {
	var blocks []components.BlockType

	if result := h.DB.Find(&blocks); result.Error != nil {
		fmt.Println(result.Error)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(blocks)
}

func (h handler) AddGenesisBlock(w http.ResponseWriter, r *http.Request) {
	block := components.CreateBlock(1, []byte{0})

	if result := h.DB.Create(&block); result.Error != nil {
		fmt.Println(result.Error)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Created")
}

func (h handler) AddBlock(w http.ResponseWriter, r *http.Request) {
	block := components.CreateBlock(h.getLength()+1, h.getPreviousHash())

	if result := h.DB.Create(&block); result.Error != nil {
		fmt.Println(result.Error)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Created")
}

func (h handler) getPreviousHash() []byte {
	var block components.BlockType

	if result := h.DB.Last(&block); result.Error != nil {
		fmt.Println(result.Error)
	}

	return block.CurrentHash
}

func (h handler) getLength() int {
	var block components.BlockType

	if result := h.DB.Last(&block); result.Error != nil {
		fmt.Println(result.Error)
	}

	return block.Id
}

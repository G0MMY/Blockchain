package Handlers

import (
	"blockchain/Controllers"
	"blockchain/Models"
	"encoding/json"
	"fmt"
	"net/http"
)

func (h Handler) IsChainValid(w http.ResponseWriter, r *http.Request) {
	blocks := h.getChain()

	if len(blocks) == 0 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Chain is null")
	} else {
		blockchain := &Models.Blockchain{Chain: blocks, Length: len(blocks)}

		if Controllers.IsChainValid(blockchain) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode("The chain is valid")
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode("the chain is invalid")
		}
	}
}

func (h Handler) GetChain(w http.ResponseWriter, r *http.Request) {
	blocks := h.getChain()

	if len(blocks) == 0 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Chain is null")
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(blocks)
	}
}

func (h Handler) GetChainLength(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.GetLength())
}

func (h Handler) getChain() []*Models.Block {
	var blocks []*Models.Block
	transactions := h.getTransactions()

	if result := h.DB.Find(&blocks); result.Error != nil {
		fmt.Println(result.Error)
	}

	i := 0
	for i < len(blocks) {
		for _, transaction := range transactions {
			if transaction.BlockID == blocks[i].ID {
				blocks[i].Transactions = append(blocks[i].Transactions, transaction)
			}
		}
		i += 1
	}
	return blocks
}

func (h Handler) GetLength() int {
	return h.GetLastBlock().ID
}

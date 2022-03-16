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
		outputs := append([]Models.Output{}, Models.Output{Amount: 10000, PublicKey: []byte("-----BEGIN RSA PUBLIC KEY-----MIICCgKCAgEAnqewLhvBAA3EF0Uob7FgFnuhBowR02TlysaXmDpncM8YohwfrmCtlUW5yBb3vmbFC1dYON+T7Kh76fk+kKyVhLn2L8X6bI4jkWesInI3PONIs8O+qLfBCtmQeFVWiKhdjmhcGU6Rxj1MZGEe3CF1D28DbOq5N6KHos2MZSmisLUHRPTmml5Au0xpPvoWy4Euoy7BoWjJnBnyUVWU7jh8vE9Hbtw3CvH5Fj7A3YrZPBCSBqAsG9eWIi3odBDTpXoHd0qjqsJfU6MRv5g7PK48j2RdjEB4I6BphuAlTomzJJbPAx2bQ4iNoqCKpWfu7y8weY3DI2yZlS/5IvSGgplRXH4NC3uuu9dB2JsT8TN5tvZtsvkyhXpmw2332oCgZEqN/dwFmT2Iwvolhy6BHNPW/OkJd5+aIxondx6RGWqrgCbqeOJ/IBFu8vdPk7n8jCZvst58iojvVx+PeabbzVlil7GbBfVcFn27XWINsNCNZGNo1NeCrC+EVwVLrvskiNvnUidW5R29IQEvZxujAOtgYo6sk4Q9aCgHqaLzFuX0LMJOKgNoXvOGIf9LZ3jpn52ZixMVETSxssRkH4WfP6bAIMOXl0MvLxjhvYdzVqhUoHmvhyrGtyzjG6A0cGzi3KErfQWQcjtkgpLi0GVfWBqHQGeCmWrENzcPIeNzE+v0DsECAwEAAQ==-----END RSA PUBLIC KEY-----")})
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
	transactions := h.GetMemPoolTransactions()
	memPoolTransactions := Controllers.FindBestMemPoolTransactions(transactions, 2)
	ids := Controllers.GetMemPoolTransactionsIds(memPoolTransactions)
	block := Controllers.CreateBlock(h.getPreviousHash(), Controllers.CreateTransactions(memPoolTransactions))

	if block.PreviousHash != nil {
		if result := h.DB.Create(block); result.Error != nil {
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

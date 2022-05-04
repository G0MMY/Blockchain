package Network

import (
	"blockchain/Models"
	"encoding/json"
	"fmt"
	"net/http"
)

func (handler *HandlerMiner) MineBlock(w http.ResponseWriter, r *http.Request) {
	var body Models.MineBlockRequest
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	defer r.Body.Close()

	fmt.Printf("Mining block id: %d with block hash: %x \n\n", body.LastIndex, body.Hash)

	Stop <- body.LastIndex
	MineBlock <- body
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Mining Block")
}

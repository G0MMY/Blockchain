package Network

import (
	"blockchain/Models"
	"encoding/json"
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

	MineBlock <- body
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Mining Block")
}

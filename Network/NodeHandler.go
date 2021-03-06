package Network

import (
	"encoding/json"
	"net/http"
)

func (handler *HandlerNode) AddNode(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	defer r.Body.Close()

	AddNode <- body["address"]
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Added successfully")
}

func (handler *HandlerNode) AddMiner(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	defer r.Body.Close()

	AddMiner <- body["miner"]
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Added successfully")
}

func (handler *HandlerNode) GetNetwork(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(handler.Node)
}

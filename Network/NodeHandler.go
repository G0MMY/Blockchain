package Network

import (
	"blockchain/Models"
	"encoding/json"
	"net/http"
)

func (handler *Handler) AddNeighbor(w http.ResponseWriter, r *http.Request) {
	body := make(map[string]string)
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	defer r.Body.Close()

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(handler.Node.UpdateNeighbors(body["address"]))
}

func (handler *Handler) UpdateNeighbor(w http.ResponseWriter, r *http.Request) {
	var body map[string][]byte
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	defer r.Body.Close()

	handler.Node.UpdateNodeNeighbor(Models.DecodeUpdateNeighbor(body["node"]))
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Updated successfully")
}

func (handler *Handler) AddNode(w http.ResponseWriter, r *http.Request) {
	var body map[string][]byte
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	defer r.Body.Close()

	handler.Node.UpdateAllNodes(Models.DecodeUpdateNeighbor(body["node"]))
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Added successfully")
}

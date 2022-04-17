package Network

import (
	"blockchain/Handlers"
	"blockchain/Models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type FullNode struct {
	Neighbors  []string
	Address    string
	Blockchain *Models.Blockchain
	LightNodes []string
}

func InitializeNode(privateKey []byte) {
	blockchain := Models.InitBlockchain(privateKey)
	handler := Handlers.New(blockchain)
	router := mux.NewRouter()

	router.HandleFunc("/lastHash", handler.GetLastHash).Methods(http.MethodGet)

	log.Println("running")
	http.ListenAndServe(":4000", router)
}

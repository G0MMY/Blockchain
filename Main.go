package main

import (
	"blockchain/database"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	//blockchain := components.InitializeBlockchain()
	//
	//blockchain.AddBlock()
	//
	//blockchain.AddTransaction("me", "me", 150, 10)
	//blockchain.AddTransaction("me", "me", 1096, 100)
	//blockchain.AddTransaction("me", "me", 12, 5)
	//blockchain.AddTransaction("me", "me", 10, 1)
	//blockchain.AddTransaction("me", "me", 1086, 180)
	//blockchain.AddTransaction("me", "me", 45, 4)
	//
	//blockchain.AddBlock()
	//
	//blockchain.AddTransaction("me", "me", 100, 10)
	//blockchain.AddTransaction("me", "me", 106, 1)
	//blockchain.AddTransaction("me", "me", 102, 50)
	//
	//blockchain.AddBlock()
	//blockchain.AddBlock()
	//
	//blockchain.AddTransaction("me", "me", 106, 1)
	//blockchain.AddTransaction("me", "me", 102, 50)

	db := database.ConnectDatabase()
	h := database.New(db)
	router := mux.NewRouter()

	router.HandleFunc("/block", h.GetBlocks).Methods(http.MethodGet)
	router.HandleFunc("/addBlock", h.AddBlock).Methods(http.MethodGet)
	router.HandleFunc("/addGenesisBlock", h.AddGenesisBlock).Methods(http.MethodGet)

	log.Println("running")
	http.ListenAndServe(":4000", router)
}

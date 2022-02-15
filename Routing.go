package main

import (
	"blockchain/database"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func initializeRoutes() {
	db := database.ConnectDatabase()
	h := database.New(db)
	router := mux.NewRouter()

	router.HandleFunc("/blockchain", h.GetChain).Methods(http.MethodGet)
	router.HandleFunc("/addBlock", h.AddBlock).Methods(http.MethodGet)
	router.HandleFunc("/addGenesisBlock", h.AddGenesisBlock).Methods(http.MethodGet)
	router.HandleFunc("/checkLastBlock", h.CheckLastBlock).Methods(http.MethodGet)
	router.HandleFunc("/isChainValid", h.IsChainValid).Methods(http.MethodGet)
	router.HandleFunc("/getChainLength", h.GetChainLength).Methods(http.MethodGet)

	log.Println("running")
	http.ListenAndServe(":4000", router)
}

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

func InitializeNode(privateKey []byte, port string) {
	blockchain := Models.InitBlockchain(privateKey)
	handler := Handlers.New(blockchain)
	router := mux.NewRouter()

	router.HandleFunc("/lastHash", handler.GetLastHash).Methods(http.MethodGet)
	router.HandleFunc("/memPoolTransactions/hash", handler.GetMemPoolTransactionsHash).Methods(http.MethodGet)
	router.HandleFunc("/merkleRoot/{blockHash}", handler.GetBlockMerkleRoot).Methods(http.MethodGet)

	router.HandleFunc("/memPoolTransactions", handler.GetMemPoolTransactions).Methods(http.MethodGet)
	router.HandleFunc("/block/{blockHash}", handler.GetBlock).Methods(http.MethodGet)
	router.HandleFunc("/transactions/{blockHash}", handler.GetBlockTransactions).Methods(http.MethodGet)
	router.HandleFunc("/balance/{address}", handler.GetPublicKeyBalance).Methods(http.MethodGet)
	router.HandleFunc("/chain", handler.GetChain).Methods(http.MethodGet)

	router.HandleFunc("/create/wallet", handler.CreateWallet).Methods(http.MethodGet)

	log.Println("running")
	http.ListenAndServe(":"+port, router)
}

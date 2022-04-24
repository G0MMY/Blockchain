package Network

import (
	"blockchain/Handlers"
	"blockchain/Models"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
)

type FullNode struct {
	Neighbors  []string
	Address    string
	Blockchain *Models.Blockchain
	LightNodes []string
}

func InitializeNode(port string, neighbors []string, lightNodes []string) {
	node := &FullNode{neighbors, port, nil, lightNodes}
	blockchain := Models.InitBlockchain(port)
	if blockchain == nil {
		node.GetBlockchain()
	}
	node.Blockchain = blockchain
	handler := Handlers.New(node.Blockchain)
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
	router.HandleFunc("/create/transaction", handler.CreateTransaction).Methods(http.MethodPost)
	router.HandleFunc("/create/block", handler.CreateBlock).Methods(http.MethodPost)

	log.Println("running")
	http.ListenAndServe(":"+port, router)
}

func (node *FullNode) GetBlockchain() {
	i := rand.Intn(len(node.Neighbors))
	resp, err := http.Get("http://localhost:" + node.Neighbors[i] + "/lastHash")
	if err != nil {
		log.Panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	bodyString := fmt.Sprintf("%s", body)

	j := 0
	for j < len(node.Neighbors) {
		if j != i {
			res, err := http.Get("http://localhost:" + node.Neighbors[j] + "/lastHash")
			if err != nil {
				log.Panic(err)
			}
			body, err = ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatalln(err)
			}
			if bodyString != fmt.Sprintf("%s", body) {
				log.Panic("There is a node with an invalid last hash")
			}
		}
		j += 1
	}

	fmt.Println(bodyString)
}

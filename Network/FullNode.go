package Network

import (
	"blockchain/Handlers"
	"blockchain/Models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
)

type FullNode struct {
	Neighbors  []string
	Address    string
	Blockchain *Models.Blockchain
}

func InitializeNode(port string, neighbors []string) {
	node := &FullNode{neighbors, port, nil}
	blockchain := Models.InitBlockchain(port)
	node.Blockchain = blockchain
	if blockchain.LastHash == nil {
		node.GetBlockchain()
	}

	handler := Handlers.New(node.Blockchain)
	router := mux.NewRouter()

	router.HandleFunc("/lastHash", handler.GetLastHash).Methods(http.MethodGet)
	router.HandleFunc("/memPoolTransactions/hash", handler.GetMemPoolTransactionsHash).Methods(http.MethodGet)
	router.HandleFunc("/merkleRoot/{blockHash}", handler.GetBlockMerkleRoot).Methods(http.MethodGet)
	router.HandleFunc("/unspentOutputs/hash", handler.GetAllUnspentOutputsHash).Methods(http.MethodGet)

	router.HandleFunc("/memPoolTransactions", handler.GetMemPoolTransactions).Methods(http.MethodGet)
	router.HandleFunc("/block/{blockHash}", handler.GetBlock).Methods(http.MethodGet)
	router.HandleFunc("/transactions/{blockHash}", handler.GetBlockTransactions).Methods(http.MethodGet)
	router.HandleFunc("/balance/{address}", handler.GetPublicKeyBalance).Methods(http.MethodGet)
	router.HandleFunc("/chain", handler.GetChain).Methods(http.MethodGet)
	router.HandleFunc("/unspentOutputs", handler.GetAllUnspentOutputs).Methods(http.MethodGet)

	router.HandleFunc("/create/wallet", handler.CreateWallet).Methods(http.MethodGet)
	router.HandleFunc("/create/transaction", handler.CreateTransaction).Methods(http.MethodPost)
	router.HandleFunc("/create/block", handler.CreateBlock).Methods(http.MethodPost)

	log.Println("running")
	http.ListenAndServe(":"+port, router)
}

func (node *FullNode) GetBlockchain() {
	lastHash := node.InitializeLastHash()
	node.InitializeBlocks(lastHash)
	node.InitializeUnspentOutputs()
	node.InitializeMemPoolTransactions()
}

func (node *FullNode) InitializeLastHash() string {
	neighbor := rand.Intn(len(node.Neighbors))

	bodyString := string(Models.ExecuteGet("http://localhost:" + node.Neighbors[neighbor] + "/lastHash"))
	bodyString = bodyString[1 : len(bodyString)-2]

	j := 0
	for j < len(node.Neighbors) {
		if j != neighbor {
			body := string(Models.ExecuteGet("http://localhost:" + node.Neighbors[j] + "/lastHash"))
			if bodyString != body[1:len(body)-2] {
				log.Panic("There is a node with an invalid last hash")
			}
		}
		j += 1
	}

	return bodyString
}

func (node *FullNode) InitializeMemPoolTransactions() {
	neighbor := rand.Intn(len(node.Neighbors))
	body := Models.ExecuteGet("http://localhost:" + string(node.Neighbors[neighbor]) + "/memPoolTransactions")
	var transactions []*Models.Transaction
	json.Unmarshal(body, &transactions)

	transactionsHash := fmt.Sprintf("%x", Models.HashTransactions(transactions))
	j := 0
	for j < len(node.Neighbors) {
		if j != neighbor {
			hashString := string(Models.ExecuteGet("http://localhost:" + string(node.Neighbors[neighbor]) + "/memPoolTransactions/hash"))

			if transactionsHash != hashString[1:len(hashString)-2] {
				log.Panic("There is a node with an invalid unspentOutput hash")
			}
		}
		j += 1
	}

	node.Blockchain.DownloadMemPool(transactions)
}

func (node *FullNode) InitializeUnspentOutputs() {
	neighbor := rand.Intn(len(node.Neighbors))
	body := Models.ExecuteGet("http://localhost:" + string(node.Neighbors[neighbor]) + "/unspentOutputs")
	unspentOutputs := make(map[string]*Models.UnspentOutput)
	json.Unmarshal(body, &unspentOutputs)

	unspentOutputHashString := fmt.Sprintf("%x", Models.HashUnspentOutputs(unspentOutputs))
	j := 0
	for j < len(node.Neighbors) {
		if j != neighbor {
			hashString := string(Models.ExecuteGet("http://localhost:" + string(node.Neighbors[neighbor]) + "/unspentOutputs/hash"))

			if unspentOutputHashString != hashString[1:len(hashString)-2] {
				log.Panic("There is a node with an invalid unspentOutput hash")
			}
		}
		j += 1
	}

	node.Blockchain.DownloadUnspentOutputs(unspentOutputs)
}

func (node *FullNode) InitializeBlocks(lastHash string) {
	neighbor := rand.Intn(len(node.Neighbors))
	body := Models.ExecuteGet("http://localhost:" + string(node.Neighbors[neighbor]) + "/block/" + lastHash)
	blockRequest := Models.BlockRequest{}
	json.Unmarshal(body, &blockRequest)

	block := blockRequest.CreateBlock()
	lastBlock := block
	node.Blockchain.PersistBlock(block)

	for true {
		if bytes.Compare(blockRequest.PreviousHash, []byte{}) == 0 {
			break
		}

		lastHash = fmt.Sprintf("%x", blockRequest.PreviousHash)
		body = Models.ExecuteGet("http://localhost:" + string(node.Neighbors[neighbor]) + "/block/" + lastHash)

		blockRequest = Models.BlockRequest{}
		json.Unmarshal(body, &blockRequest)
		block = blockRequest.CreateBlock()

		if bytes.Compare(lastBlock.PreviousHash, block.Hash()) != 0 {
			log.Panic("Wrong block")
		}
		lastBlock = block
		node.Blockchain.PersistBlock(block)
	}

	node.Blockchain.ValidateBlockchain()
}

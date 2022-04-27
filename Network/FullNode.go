package Network

import (
	"blockchain/Models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
)

var (
	minNeighbors = 5
)

type FullNode struct {
	Neighbors  []string
	AllNodes   []Models.OtherFullNode
	Address    string
	Blockchain *Models.Blockchain
}

type numberNeighbor struct {
	address string
	number  int
}

func InitializeNode(port string, allNodes []Models.OtherFullNode) {
	node := &FullNode{nil, allNodes, port, nil}
	node.getNeighbors()
	node.addToNetwork()

	//node.Blockchain = Models.InitTestBlockchain(Models.CreateWallet().PrivateKey)
	node.Blockchain = Models.InitBlockchain(port)
	if node.Blockchain.LastHash == nil {
		node.getBlockchain()
	}

	handler := New(node)
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

	router.HandleFunc("/update/neighbor", handler.UpdateNeighbor).Methods(http.MethodPost)

	router.HandleFunc("/add/neighbor", handler.AddNeighbor).Methods(http.MethodPost)
	router.HandleFunc("/add/node", handler.AddNode).Methods(http.MethodPost)

	log.Println("running")
	http.ListenAndServe(":"+port, router)
}

func (node *FullNode) getNeighbors() {
	if node.AllNodes == nil {
		return
	}

	numberNeighborsPerNode := make(map[string]int)

	for _, otherNode := range node.AllNodes {
		if len(otherNode.Neighbors) == 0 {
			numberNeighborsPerNode[otherNode.Address] = 0
		}
		for _, address := range otherNode.Neighbors {
			if val, ok := numberNeighborsPerNode[address]; ok {
				numberNeighborsPerNode[address] = val + 1
			} else {
				numberNeighborsPerNode[address] = 1
			}
		}
	}

	var otherNodeArray []numberNeighbor
	for key, value := range numberNeighborsPerNode {
		if len(otherNodeArray) < minNeighbors {
			otherNodeArray = append(otherNodeArray, numberNeighbor{key, value})
		} else {
			for i, otherNode := range otherNodeArray {
				if otherNode.number > value {
					otherNodeArray = append(otherNodeArray[:i+1], otherNodeArray[i:]...)
					otherNodeArray[i] = numberNeighbor{key, value}
				}
			}
		}
	}

	for _, otherNode := range otherNodeArray {
		node.UpdateNeighbors(otherNode.address)
	}
}

func (node *FullNode) UpdateAllNodes(newNode Models.UpdateNetwork) {
	neighbors := node.Neighbors
	for _, received := range newNode.Received {
		if received == node.Address {
			return
		} else {
			for i, neighbor := range neighbors {
				if neighbor == received {
					neighbors = append(neighbors[:i], neighbors[i+1:]...)
					break
				}
			}
		}
	}
	for _, otherNode := range node.AllNodes {
		if otherNode.Address == newNode.Node.Address {
			return
		}
	}

	node.AllNodes = append(node.AllNodes, newNode.Node)
	newNode.Received = append(newNode.Received, node.Address)

	for _, neighbor := range neighbors {
		postBody, _ := json.Marshal(map[string][]byte{
			"node": newNode.Encode(),
		})
		responseBody := bytes.NewBuffer(postBody)
		go Models.ExecutePost("http://localhost:"+neighbor+"/add/node", responseBody)
	}
}

func (node *FullNode) UpdateNodeNeighbor(update Models.UpdateNetwork) {
	neighbors := node.Neighbors
	for _, received := range update.Received {
		if received == node.Address {
			return
		} else {
			for i, neighbor := range neighbors {
				if neighbor == received {
					neighbors = append(neighbors[:i], neighbors[i+1:]...)
					break
				}
			}
		}
	}
	update.Received = append(update.Received, node.Address)

	for i, otherFullNode := range node.AllNodes {
		if update.Node.Address == otherFullNode.Address {
			node.AllNodes[i] = update.Node
			break
		}
	}

	for _, neighbor := range neighbors {
		postBody, _ := json.Marshal(map[string][]byte{
			"node": update.Encode(),
		})
		responseBody := bytes.NewBuffer(postBody)
		go Models.ExecutePost("http://localhost:"+neighbor+"/update/neighbor", responseBody)
	}
}

func (node *FullNode) UpdateNeighbors(address string) []string {
	for _, neighbor := range node.Neighbors {
		if neighbor == address || address == node.Address {
			return node.Neighbors
		}
	}

	neighbors := append(node.Neighbors, address)
	for _, neighbor := range node.Neighbors {
		postBody, _ := json.Marshal(map[string][]byte{
			"node": Models.UpdateNetwork{[]string{node.Address}, Models.OtherFullNode{node.Address, neighbors}}.Encode(),
		})
		responseBody := bytes.NewBuffer(postBody)
		go Models.ExecutePost("http://localhost:"+neighbor+"/update/neighbor", responseBody)
	}

	node.Neighbors = neighbors
	for i, otherNode := range node.AllNodes {
		if node.Address == otherNode.Address {
			node.AllNodes[i].Neighbors = neighbors
			break
		}
	}

	return node.Neighbors
}

func (node *FullNode) addToNetwork() {
	if node.AllNodes == nil {
		node.AllNodes = append(node.AllNodes, Models.OtherFullNode{node.Address, nil})
		return
	}

	var otherNodes []Models.OtherFullNode
	for _, otherNode := range node.AllNodes {
		if otherNode.Address != node.Address {
			if len(otherNodes) == 0 {
				otherNodes = append(otherNodes, otherNode)
			} else {
				for i, otherFullNode := range otherNodes {
					if len(otherFullNode.Neighbors) > len(otherNode.Neighbors) {
						otherNodes = append(otherNodes[:i+1], otherNodes[i:]...)
						otherNodes[i] = otherNode
						break
					}
				}
			}
		}
	}
	newNode := Models.UpdateNetwork{[]string{node.Address}, Models.OtherFullNode{node.Address, node.Neighbors}}
	node.AllNodes = append(node.AllNodes, Models.OtherFullNode{node.Address, node.Neighbors})

	i := 0
	for i < minNeighbors {
		if i >= len(otherNodes) {
			return
		}
		for j, otherNode := range node.AllNodes {
			if otherNode.Address == otherNodes[i].Address {
				node.AllNodes[j] = Models.OtherFullNode{otherNode.Address, append(otherNode.Neighbors, node.Address)}
			}
		}
		postBody, _ := json.Marshal(map[string]string{
			"address": node.Address,
		})
		responseBody := bytes.NewBuffer(postBody)
		Models.ExecutePost("http://localhost:"+otherNodes[i].Address+"/add/neighbor", responseBody)

		postBody, _ = json.Marshal(map[string][]byte{
			"node": newNode.Encode(),
		})
		responseBody = bytes.NewBuffer(postBody)
		Models.ExecutePost("http://localhost:"+otherNodes[i].Address+"/add/node", responseBody)

		i += 1
	}
}

func (node *FullNode) getBlockchain() {
	lastHash := node.initializeLastHash()
	node.initializeBlocks(lastHash)
	node.initializeUnspentOutputs()
	node.initializeMemPoolTransactions()
}

func (node *FullNode) initializeLastHash() string {
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

func (node *FullNode) initializeMemPoolTransactions() {
	neighbor := rand.Intn(len(node.Neighbors))
	body := Models.ExecuteGet("http://localhost:" + node.Neighbors[neighbor] + "/memPoolTransactions")
	var transactions []*Models.Transaction
	json.Unmarshal(body, &transactions)

	transactionsHash := fmt.Sprintf("%x", Models.HashTransactions(transactions))
	j := 0
	for j < len(node.Neighbors) {
		if j != neighbor {
			hashString := string(Models.ExecuteGet("http://localhost:" + node.Neighbors[neighbor] + "/memPoolTransactions/hash"))

			if transactionsHash != hashString[1:len(hashString)-2] {
				log.Panic("There is a node with an invalid unspentOutput hash")
			}
		}
		j += 1
	}

	node.Blockchain.DownloadMemPool(transactions)
}

func (node *FullNode) initializeUnspentOutputs() {
	neighbor := rand.Intn(len(node.Neighbors))
	body := Models.ExecuteGet("http://localhost:" + node.Neighbors[neighbor] + "/unspentOutputs")
	unspentOutputs := make(map[string]*Models.UnspentOutput)
	json.Unmarshal(body, &unspentOutputs)

	unspentOutputHashString := fmt.Sprintf("%x", Models.HashUnspentOutputs(unspentOutputs))
	j := 0
	for j < len(node.Neighbors) {
		if j != neighbor {
			hashString := string(Models.ExecuteGet("http://localhost:" + node.Neighbors[neighbor] + "/unspentOutputs/hash"))

			if unspentOutputHashString != hashString[1:len(hashString)-2] {
				log.Panic("There is a node with an invalid unspentOutput hash")
			}
		}
		j += 1
	}

	node.Blockchain.DownloadUnspentOutputs(unspentOutputs)
}

func (node *FullNode) initializeBlocks(lastHash string) {
	neighbor := rand.Intn(len(node.Neighbors))
	body := Models.ExecuteGet("http://localhost:" + node.Neighbors[neighbor] + "/block/" + lastHash)
	blockRequest := Models.BlockRequest{}
	json.Unmarshal(body, &blockRequest)

	block := blockRequest.CreateBlock()
	if !block.Validate() {
		log.Panic("Invalid block")
	}
	lastBlock := block
	node.Blockchain.PersistBlock(block)

	for true {
		if bytes.Compare(blockRequest.PreviousHash, []byte{}) == 0 {
			break
		}

		lastHash = fmt.Sprintf("%x", blockRequest.PreviousHash)
		body = Models.ExecuteGet("http://localhost:" + node.Neighbors[neighbor] + "/block/" + lastHash)

		blockRequest = Models.BlockRequest{}
		json.Unmarshal(body, &blockRequest)
		block = blockRequest.CreateBlock()
		if !block.Validate() {
			log.Panic("Invalid block")
		}

		if bytes.Compare(lastBlock.PreviousHash, block.Hash()) != 0 {
			log.Panic("Wrong block")
		}
		lastBlock = block
		node.Blockchain.PersistBlock(block)
	}

	node.Blockchain.ValidateBlockchain()
}

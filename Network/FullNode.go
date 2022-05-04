package Network

import (
	"blockchain/Models"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
)

var (
	minChecks         = 5
	AddNode           chan string
	CreateBlock       chan Models.BlockRequest
	AddBlock          chan Models.BlockRequest
	CreateTransaction chan Models.TransactionRequest
	AddTransaction    chan Models.CreateTransactionRequest
	AddMiner          chan string
)

type FullNode struct {
	AllNodes   []string
	Address    string
	Miners     []string
	Blockchain *Models.Blockchain
}

func InitializeNode(port, neighbor string) {
	allNodes := getNetwork(neighbor)
	for i, node := range allNodes {
		if node == port {
			allNodes = append(allNodes[:i], allNodes[i+1:]...)
		}
	}
	node := &FullNode{allNodes, port, []string{}, nil}

	node.Blockchain = Models.InitBlockchain(port)
	if node.Blockchain == nil {
		return
	} else if node.Blockchain.LastHash == nil {
		if !node.getBlockchain(nil) {
			return
		}
	} else if allNodes != nil {
		lastHash := node.initializeLastHash()
		if lastHash != fmt.Sprintf("%x", node.Blockchain.LastHash) {
			if !node.updateBlockchain(node.Blockchain.LastHash) {
				return
			}
		}
	}
	node.initializeAllNodes()

	handler := NewNode(node)
	router := mux.NewRouter()

	router.HandleFunc("/lastHash", handler.GetLastHash).Methods(http.MethodGet)
	router.HandleFunc("/memPoolTransactions/hash", handler.GetMemPoolTransactionsHash).Methods(http.MethodGet)
	router.HandleFunc("/merkleRoot/{blockHash}", handler.GetBlockMerkleRoot).Methods(http.MethodGet)
	router.HandleFunc("/unspentOutputs/hash", handler.GetAllUnspentOutputsHash).Methods(http.MethodGet)
	router.HandleFunc("/network", handler.GetNetwork).Methods(http.MethodGet)

	router.HandleFunc("/memPoolTransactions", handler.GetMemPoolTransactions).Methods(http.MethodGet)
	router.HandleFunc("/block/{blockHash}", handler.GetBlock).Methods(http.MethodGet)
	router.HandleFunc("/lastBLock", handler.GetLastBlock).Methods(http.MethodGet)
	router.HandleFunc("/transactions/{blockHash}", handler.GetBlockTransactions).Methods(http.MethodGet)
	router.HandleFunc("/balance/{address}", handler.GetPublicKeyBalance).Methods(http.MethodGet)
	router.HandleFunc("/chain", handler.GetChain).Methods(http.MethodGet)
	router.HandleFunc("/unspentOutputs", handler.GetAllUnspentOutputs).Methods(http.MethodGet)

	router.HandleFunc("/create/wallet", handler.CreateWallet).Methods(http.MethodGet)
	router.HandleFunc("/create/transaction", handler.CreateTransaction).Methods(http.MethodPost)
	router.HandleFunc("/create/block", handler.CreateBlock).Methods(http.MethodPost)

	router.HandleFunc("/add/node", handler.AddNode).Methods(http.MethodPost)
	router.HandleFunc("/add/block", handler.AddBlock).Methods(http.MethodPost)
	router.HandleFunc("/add/transaction", handler.AddTransaction).Methods(http.MethodPost)
	router.HandleFunc("/add/miner", handler.AddMiner).Methods(http.MethodPost)

	AddNode = make(chan string, 100)
	CreateBlock = make(chan Models.BlockRequest, 100)
	AddBlock = make(chan Models.BlockRequest, 100)
	AddTransaction = make(chan Models.CreateTransactionRequest, 100)
	CreateTransaction = make(chan Models.TransactionRequest, 100)
	AddMiner = make(chan string, 100)

	go addNodeWorker(AddNode, node)
	go createBlockWorker(CreateBlock, node)
	go addBLockWorker(AddBlock, node)
	go addTransactionWorker(AddTransaction, node)
	go createTransactionWorker(CreateTransaction, node)
	go addMinerWorker(AddMiner, node)

	log.Println("running")
	http.ListenAndServe(":"+port, router)
}

func addMinerWorker(addMiner <-chan string, node *FullNode) {
	for miner := range addMiner {
		node.addMiner(miner)
	}
}

func createBlockWorker(createBlock <-chan Models.BlockRequest, node *FullNode) {
	for blockRequest := range createBlock {
		node.createBlock(blockRequest)
	}
}

func addBLockWorker(addBlock <-chan Models.BlockRequest, node *FullNode) {
	for block := range addBlock {
		node.addBlock(block)
	}
}

func addNodeWorker(addNode <-chan string, node *FullNode) {
	for update := range addNode {
		node.addNode(update)
	}
}

func addTransactionWorker(addTransaction <-chan Models.CreateTransactionRequest, node *FullNode) {
	for transaction := range addTransaction {
		node.addTransaction(transaction)
	}
}

func createTransactionWorker(createTransaction <-chan Models.TransactionRequest, node *FullNode) {
	for transaction := range createTransaction {
		node.createTransaction(transaction)
	}
}

func (node *FullNode) updateBlockchain(currentLastHash []byte) bool {
	node.Blockchain.LastHash = nil
	if !node.Blockchain.Delete(Models.MemPoolPrefix) || !node.Blockchain.Delete(Models.UnspentOutputPrefix) {
		return false
	} else if !node.getBlockchain(currentLastHash) {
		return false
	}

	return true
}

func (node *FullNode) addMiner(miner string) {
	memPool := node.Blockchain.GetMemPoolTransactions()
	lastBlock := node.Blockchain.GetLastBlock()
	hash := lastBlock.Hash()
	if hash == nil {
		return
	}
	sendBlockToMiner(memPool, lastBlock, hash, miner)

	for _, address := range node.Miners {
		if address == miner {
			return
		}
	}

	node.Miners = append(node.Miners, miner)
}

func getNetwork(neighbor string) []string {
	if neighbor == "" {
		return nil
	}

	var node FullNode
	byteNetwork := Models.ExecuteGet("http://localhost:" + neighbor + "/network")
	json.Unmarshal(byteNetwork, &node)

	return node.AllNodes
}

func (node *FullNode) createTransaction(transactionRequest Models.TransactionRequest) {
	transaction := transactionRequest.CreateTransaction()
	if !transaction.ValidateTransaction(true) {
		return
	}
	node.Blockchain.PersistTransaction(transaction)
}

func (node *FullNode) addTransaction(transactionRequest Models.CreateTransactionRequest) {
	priv, err := hex.DecodeString(transactionRequest.PrivateKey)
	if err != nil {
		log.Println(err)
		return
	}

	to, err := hex.DecodeString(transactionRequest.To)
	if err != nil {
		log.Println(err)
		return
	}

	transaction := node.Blockchain.CreateTransaction(priv, to, transactionRequest.Amount, transactionRequest.Fee, transactionRequest.Timestamp)
	if transaction != nil {
		for _, otherNode := range node.AllNodes {
			if otherNode != node.Address {
				byteTransaction := transaction.EncodeTransaction()
				if byteTransaction == nil {
					return
				}

				body := bytes.NewBuffer(byteTransaction)
				go Models.ExecutePost("http://localhost:"+otherNode+"/create/transaction", body)
			}
		}
	} else {
		log.Println("Bad Transaction")
	}
}

func (node *FullNode) sendBlockToMiners() {
	if len(node.Miners) == 0 {
		return
	}
	memPool := node.Blockchain.GetMemPoolTransactions()
	lastBlock := node.Blockchain.GetLastBlock()
	hash := lastBlock.Hash()
	if hash == nil {
		return
	}

	for _, miner := range node.Miners {
		sendBlockToMiner(memPool, lastBlock, hash, miner)
	}
}

func sendBlockToMiner(memPool []*Models.Transaction, lastBlock *Models.Block, hash []byte, miner string) {
	byteBody, err := json.Marshal(Models.MineBlockRequest{lastBlock.Index, hash, memPool})
	if err != nil {
		log.Println(err)
		return
	}
	body := bytes.NewBuffer(byteBody)
	Models.ExecutePost("http://localhost:"+miner+"/mine/block", body)
}

func (node *FullNode) createBlock(blockRequest Models.BlockRequest) {
	block := blockRequest.CreateBlock()
	if block == nil {
		return
	}
	node.Blockchain.CreateBlock(block)
	node.sendBlockToMiners()
}

func (node *FullNode) addBlock(blockRequest Models.BlockRequest) {
	block := blockRequest.CreateBlock()
	if block == nil {
		return
	}
	_, err := node.Blockchain.CreateBlock(block)

	if err == "" {
		for _, otherNode := range node.AllNodes {
			if otherNode != node.Address {
				body := bytes.NewBuffer(block.EncodeBlock())
				go Models.ExecutePost("http://localhost:"+otherNode+"/create/block", body)
			}
		}
		node.sendBlockToMiners()
	} else {
		log.Println(err)
	}
}

func (node *FullNode) addNode(address string) {
	for _, currentNode := range node.AllNodes {
		if currentNode == address {
			return
		}
	}

	node.AllNodes = append(node.AllNodes, address)
}

func (node *FullNode) initializeAllNodes() {
	for _, currentNode := range node.AllNodes {
		if currentNode == node.Address {
			return
		}
	}

	for _, currentNode := range node.AllNodes {
		byteBody, err := json.Marshal(map[string]string{
			"address": node.Address,
		})

		if err != nil {
			log.Println(err)
			return
		}

		body := bytes.NewBuffer(byteBody)
		go Models.ExecutePost("http://localhost:"+currentNode+"/add/node", body)
	}
	node.AllNodes = append(node.AllNodes, node.Address)
}

func (node *FullNode) getBlockchain(currentLastHash []byte) bool {
	lastHash := node.initializeLastHash()
	if lastHash == "" {
		return false
	}
	if !node.initializeBlocks(lastHash, currentLastHash) {
		return false
	}
	if !node.initializeUnspentOutputs() {
		return false
	}
	if !node.initializeMemPoolTransactions() {
		return false
	}

	return true
}

func (node *FullNode) initializeLastHash() string {
	neighbor := rand.Intn(len(node.AllNodes))

	lastHash := Models.ExecuteGet("http://localhost:" + node.AllNodes[neighbor] + "/lastHash")
	if lastHash == nil {
		return ""
	}
	bodyString := string(lastHash)
	bodyString = bodyString[1 : len(bodyString)-2]

	j := 0
	for j < minChecks && j < len(node.AllNodes) {
		if j != neighbor {
			lastHash = Models.ExecuteGet("http://localhost:" + node.AllNodes[j] + "/lastHash")
			if lastHash == nil {
				return ""
			}
			body := string(lastHash)
			if bodyString != body[1:len(body)-2] {
				log.Println("There is a node with an invalid last hash")
				return ""
			}
		}
		j += 1
	}

	return bodyString
}

func (node *FullNode) initializeMemPoolTransactions() bool {
	neighbor := rand.Intn(len(node.AllNodes))
	body := Models.ExecuteGet("http://localhost:" + node.AllNodes[neighbor] + "/memPoolTransactions")
	if body == nil {
		return false
	}
	var transactions []*Models.Transaction
	json.Unmarshal(body, &transactions)

	hash := Models.HashTransactions(transactions)
	if hash == nil {
		return false
	}
	transactionsHash := fmt.Sprintf("%x", hash)
	j := 0
	for j < minChecks && j < len(node.AllNodes) {
		if j != neighbor {
			hash := Models.ExecuteGet("http://localhost:" + node.AllNodes[neighbor] + "/memPoolTransactions/hash")
			if hash == nil {
				return false
			}
			hashString := string(hash)

			if transactionsHash != hashString[1:len(hashString)-2] {
				log.Println("There is a node with an invalid unspentOutput hash")
				return false
			}
		}
		j += 1
	}
	node.Blockchain.DownloadMemPool(transactions)

	return true
}

func (node *FullNode) initializeUnspentOutputs() bool {
	neighbor := rand.Intn(len(node.AllNodes))
	body := Models.ExecuteGet("http://localhost:" + node.AllNodes[neighbor] + "/unspentOutputs")
	if body == nil {
		return false
	}
	unspentOutputs := make(map[string]*Models.UnspentOutput)
	json.Unmarshal(body, &unspentOutputs)

	unspentOutputHashString := fmt.Sprintf("%x", Models.HashUnspentOutputs(unspentOutputs))
	j := 0
	for j < minChecks && j < len(node.AllNodes) {
		if j != neighbor {
			hash := Models.ExecuteGet("http://localhost:" + node.AllNodes[neighbor] + "/unspentOutputs/hash")
			if hash == nil {
				return false
			}
			hashString := string(hash)

			if unspentOutputHashString != hashString[1:len(hashString)-2] {
				log.Println("There is a node with an invalid unspentOutput hash")
				return false
			}
		}
		j += 1
	}
	node.Blockchain.DownloadUnspentOutputs(unspentOutputs)

	return true
}

func (node *FullNode) initializeBlocks(lastHash string, currentLastHash []byte) bool {
	neighbor := rand.Intn(len(node.AllNodes))
	body := Models.ExecuteGet("http://localhost:" + node.AllNodes[neighbor] + "/block/" + lastHash)
	if body == nil {
		return false
	}
	blockRequest := Models.BlockRequest{}
	json.Unmarshal(body, &blockRequest)

	block := blockRequest.CreateBlock()
	if block == nil {
		return false
	}
	if !block.Validate() {
		log.Println("Invalid block")
		return false
	}
	lastBlock := block
	node.Blockchain.PersistBlock(block)

	var compare []byte
	if currentLastHash == nil {
		compare = []byte{}
	} else {
		compare = currentLastHash
	}

	for true {
		if bytes.Compare(blockRequest.PreviousHash, compare) == 0 {
			break
		}

		lastHash = fmt.Sprintf("%x", blockRequest.PreviousHash)
		body = Models.ExecuteGet("http://localhost:" + node.AllNodes[neighbor] + "/block/" + lastHash)
		if body == nil {
			return false
		}

		blockRequest = Models.BlockRequest{}
		json.Unmarshal(body, &blockRequest)
		block = blockRequest.CreateBlock()
		if block == nil {
			return false
		}
		if !block.Validate() {
			log.Println("Invalid block")
			return false
		}

		hash := block.Hash()
		if hash == nil {
			return false
		} else if bytes.Compare(lastBlock.PreviousHash, hash) != 0 {
			log.Println("Wrong block")
			return false
		}

		lastBlock = block
		node.Blockchain.PersistBlock(block)
	}

	return node.Blockchain.ValidateBlockchain()
}

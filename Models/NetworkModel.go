package Models

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

type WalletResponse struct {
	PrivateKey string
	PublicKey  string
	Address    string
}

type CreateTransactionRequest struct {
	PrivateKey string
	To         string
	Amount     int
	Fee        int
	Timestamp  int64
}

type BlockRequest struct {
	Index        int
	Nonce        int
	Timestamp    int64
	MerkleRoot   []byte
	PreviousHash []byte
	Transactions []TransactionRequest
}

type TransactionRequest struct {
	Inputs    []InputRequest
	Outputs   []Output
	Timestamp int64
	Fee       int
}

type InputRequest struct {
	Output    Output
	Signature []byte
	PublicKey []byte
}

type UnspentOutputsRequest struct {
	Outputs []Output
}

func (request UnspentOutputsRequest) CreateUnspentOutput() *UnspentOutput {
	var outputs []*Output

	for _, output := range request.Outputs {
		outputs = append(outputs, &output)
	}

	return &UnspentOutput{outputs}
}

func (unspentOutput *UnspentOutput) CreateUnspentOutputRequest() UnspentOutputsRequest {
	var outputs []Output

	for _, output := range unspentOutput.Outputs {
		outputs = append(outputs, *output)
	}

	return UnspentOutputsRequest{outputs}
}

func (block *Block) CreateBlockRequest() BlockRequest {
	return BlockRequest{block.Index, block.Nonce, block.Timestamp, block.MerkleRoot, block.PreviousHash, CreateTransactionsRequest(block.Transactions)}
}

func (blockRequest BlockRequest) CreateBlock() *Block {
	var transactions []*Transaction

	for _, transactionRequest := range blockRequest.Transactions {
		transactions = append(transactions, transactionRequest.CreateTransaction())
	}
	tree := CreateTree(transactions)
	if tree == nil {
		return nil
	}

	if bytes.Compare(tree.RootNode.Data, blockRequest.MerkleRoot) != 0 {
		log.Println("The tree or the merkle root is not valid")
	}

	return &Block{blockRequest.Index, blockRequest.Nonce, blockRequest.Timestamp, blockRequest.MerkleRoot, blockRequest.PreviousHash, transactions, tree}
}

func (transactionRequest TransactionRequest) CreateTransaction() *Transaction {
	var inputs []*Input
	for _, input := range transactionRequest.Inputs {
		inputs = append(inputs, &Input{&input.Output, input.Signature, input.PublicKey})
	}

	var outputs []*Output
	for _, output := range transactionRequest.Outputs {
		outputs = append(outputs, &Output{output.PublicKeyHash, output.Amount})
	}

	return &Transaction{inputs, outputs, transactionRequest.Timestamp, transactionRequest.Fee}
}

func (transaction *Transaction) CreateTransactionRequest() TransactionRequest {
	var inputs []InputRequest
	for _, input := range transaction.Inputs {
		inputs = append(inputs, InputRequest{*input.Output, input.Signature, input.PublicKey})
	}

	var outputs []Output
	for _, output := range transaction.Outputs {
		outputs = append(outputs, *output)
	}

	return TransactionRequest{inputs, outputs, transaction.Timestamp, transaction.Fee}
}

func CreateTransactionsRequest(transactions []*Transaction) []TransactionRequest {
	var transactionsRequest []TransactionRequest
	for _, transaction := range transactions {
		transactionsRequest = append(transactionsRequest, transaction.CreateTransactionRequest())
	}

	return transactionsRequest
}

func ExecutePost(url string, responseBody *bytes.Buffer) []byte {
	resp, err := http.Post(url, "application/json", responseBody)

	if err != nil {
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	} else if resp.StatusCode != 200 {
		log.Println(string(body))
		return nil
	}

	return body
}

func ExecuteGet(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	} else if resp.StatusCode != 200 {
		log.Println(string(body))
		return nil
	}

	return body
}

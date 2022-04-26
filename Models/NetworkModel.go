package Models

import (
	"bytes"
	"encoding/hex"
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

type CreateBlockResponse struct {
	Index        int
	Nonce        int
	Timestamp    int64
	MerkleRoot   string
	PreviousHash string
	Transactions []string
	MerkleTree   string
	PrivateKey   string
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

	if bytes.Compare(tree.RootNode.Data, blockRequest.MerkleRoot) != 0 {
		log.Panic("The tree or the merkle root is not valid")
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

func CreateBlockToBlock(createBlock *CreateBlockResponse) *Block {
	var transactions []*Transaction

	for _, stringHash := range createBlock.Transactions {
		hash, err := hex.DecodeString(stringHash)
		if err != nil {
			return nil
		}
		transactions = append(transactions, DecodeTransaction(hash))
	}

	merkleRoot, err := hex.DecodeString(createBlock.MerkleRoot)
	if err != nil {
		return nil
	}

	previousHash, err := hex.DecodeString(createBlock.PreviousHash)
	if err != nil {
		return nil
	}

	tree, err := hex.DecodeString(createBlock.MerkleTree)
	if err != nil {
		return nil
	}

	return &Block{createBlock.Index, createBlock.Nonce, createBlock.Timestamp, merkleRoot, previousHash, transactions, DecodeTree(tree)}
}

func ExecuteGet(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body
}

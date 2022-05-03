package Network

import (
	"blockchain/Models"
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"math/big"
	"net/http"
	"time"
)

var (
	MineBlock chan Models.MineBlockRequest
	Stop      chan int
)

type Miner struct {
	Address      string
	FullNode     string
	privateKey   []byte
	stop         bool
	CurrentBlock Models.MineBlockRequest
}

func InitializeMiner(port string, node string, privateKey []byte) {
	miner := &Miner{port, node, privateKey, false, Models.MineBlockRequest{}}

	if !miner.AddToNode(node) {
		log.Println("Bad node")
		return
	}
	handler := NewMiner(miner)
	router := mux.NewRouter()

	router.HandleFunc("/mine/block", handler.MineBlock).Methods(http.MethodPost)

	MineBlock = make(chan Models.MineBlockRequest, 100)
	Stop = make(chan int, 100)

	go mineBlockWorker(MineBlock, miner)

	log.Println("running")
	http.ListenAndServe(":"+port, router)
}

func mineBlockWorker(mineBlock <-chan Models.MineBlockRequest, miner *Miner) {
	for {
		select {
		case block := <-mineBlock:
			miner.start(block)
		}
	}
}

func (miner *Miner) start(mineBlock Models.MineBlockRequest) {
	if mineBlock.Hash != nil {
		block := Models.CreateBlockMiner(miner.privateKey, mineBlock.LastIndex+1, mineBlock.Hash, mineBlock.MemPoolTransactions)
		if miner.mine(block) {
			body := bytes.NewBuffer(block.EncodeBlock())
			Models.ExecutePost("http://localhost:"+miner.FullNode+"/add/block", body)
		}
	}
}

func (miner *Miner) mine(block *Models.Block) bool {
	var intHash big.Int
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Models.Difficulty))

	time.Sleep(5 * time.Second)
	for {
		select {
		case index := <-Stop:
			if index != block.Index-1 {
				return false
			}
		default:
			block.Nonce += 1
			hash := block.Hash()
			if hash == nil {
				block.Nonce = -1
				return false
			}

			intHash.SetBytes(hash)
			if intHash.Cmp(target) == -1 {
				return true
			}
		}
	}
}

func (miner *Miner) AddToNode(node string) bool {
	byteBody, err := json.Marshal(map[string]string{
		"miner": miner.Address,
	})

	if err != nil {
		log.Println(err)
		return false
	}

	body := bytes.NewBuffer(byteBody)
	Models.ExecutePost("http://localhost:"+node+"/add/miner", body)

	return true
}

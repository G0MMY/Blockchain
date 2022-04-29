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
)

type Miner struct {
	Address      string
	FullNode     string
	running      bool
	privateKey   []byte
	CurrentBlock Models.MineBlockRequest
}

func InitializeMiner(port string, node string, privateKey []byte) {
	miner := &Miner{port, node, false, privateKey, Models.MineBlockRequest{}}

	if !miner.AddToNode(node) {
		log.Println("Bad node")
		return
	}

	handler := NewMiner(miner)
	router := mux.NewRouter()

	router.HandleFunc("/mine/block", handler.MineBlock).Methods(http.MethodPost)

	MineBlock = make(chan Models.MineBlockRequest, 100)

	go mineBlockWorker(MineBlock, miner)

	log.Println("running")
	http.ListenAndServe(":"+port, router)
}

func mineBlockWorker(mineBlock <-chan Models.MineBlockRequest, miner *Miner) {
	for mineBlockRequest := range mineBlock {
		miner.Stop()
		miner.CurrentBlock = mineBlockRequest
		go miner.Start()
	}
}

func (miner *Miner) Start() {
	if miner.CurrentBlock.Hash != nil {
		miner.running = true
		block := Models.CreateBlockMiner(miner.privateKey, miner.CurrentBlock.LastIndex+1, miner.CurrentBlock.Hash, miner.CurrentBlock.MemPoolTransactions)
		if miner.Mine(block) {
			body := bytes.NewBuffer(block.EncodeBlock())
			Models.ExecutePost("http://localhost:"+miner.FullNode+"/add/block", body)
		}
	}
	miner.running = false
}

func (miner *Miner) Mine(block *Models.Block) bool {
	if miner.running {
		var intHash big.Int
		target := big.NewInt(1)
		target.Lsh(target, uint(256-Models.Difficulty))

		time.Sleep(10 * time.Second)
		for true {
			if !miner.running {
				return false
			}
			block.Nonce += 1
			hash := block.Hash()
			if hash == nil {
				block.Nonce = -1
				return false
			}

			intHash.SetBytes(hash)
			if intHash.Cmp(target) == -1 {
				break
			}
		}
	}

	return true
}

func (miner *Miner) Stop() {
	miner.running = false
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

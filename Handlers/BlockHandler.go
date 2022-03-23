package Handlers

import (
	"blockchain/Controllers"
	"blockchain/Models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (h Handler) AddGenesisBlock(w http.ResponseWriter, r *http.Request) {
	if h.GetLength() > 0 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Can't add Genesis Block on top of existing blocks")
	} else {
		outputs := append([]Models.Output{}, Models.Output{Amount: 10000, PublicKey: []byte("MIICCgKCAgEAuxhXQdgiCIKCtFovJ7QNBXWCG8qjRQWY55/Ci+DRTnb23EXkz3bQFp4vJpU/CVh6UORm7Gep1MAXGR0WLH4q4joLuPrxFYTUnUyMy++Fy7x7eWixioCFR7ySv6qmb4hYunVtbj0MoHBDkg3MDlwoaqMT24Imgd+Be1MBtmBNVW9uF2kCtoKnGKyxhASAEpQ5EPDYIp57IZaFjKDa3dlNFIPqKeoFwRoZ3qlBUV3bY28GVu+4fyuxisnPOSpm3uGpPgsqabdmcoce0LgNJYOiXbGjSY2RrLr2j1eBFx8aAJGdkEgZY21UeDi0rMFuLEzlXcqbBksw4sTA6GFr1ngKBKZ6PYiwNgmYndw8kGowKfWIuOzR8OFZV+KMn5mmcUpk14u2gT1b+BQG/T+4bXz/71ojEOTJESjfZ8CV9kBAx/dM/5GN3qjSvtA0eK8oUUmlgbWuIeC4SRMwI7WD3RxG0fi23SgrzCwM20m/qOlx6zo7oJ2BvgFlqPD3FDwVtofpEyxYha4EbUNNuLbi0zcLwn1GKv6vLI0AVv4pJXXRIg3WKjm3kRGWx3IYzqTtFGbtN2BQjbIm2LEWGzhd6CvC/WaP+08kPAUO6hrWe5DbcIlEPExKCTxNzsO1dfOwAw6KSlFGUSUTeC3dm1eavxpEQg35ekGSSITYGVeq/0lDyHkCAwEAAQ==")})
		transactions := append([]Models.Transaction{}, Models.Transaction{Outputs: outputs, Timestamp: time.Now().Unix()})
		block := Controllers.CreateBlock([]byte{0}, transactions)

		if result := h.DB.Create(&block); result.Error != nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(result.Error)
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(block)
		}
	}
}

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

func (h Handler) AddBlock(w http.ResponseWriter, r *http.Request) {
	transactions := h.GetMemPoolTransactions()
	memPoolTransactions := Controllers.FindBestMemPoolTransactions(transactions, 2)
	ids := Controllers.GetMemPoolTransactionsIds(memPoolTransactions)
	block := Controllers.CreateBlock(h.getPreviousHash(), Controllers.CreateTransactions(memPoolTransactions))

	if block.PreviousHash != nil {
		if result := h.DB.Create(block); result.Error != nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(result.Error)
		} else {
			if h.DeleteMemPoolTransactions(ids) {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(block)
			} else {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("There was an error with the transactions")
			}
		}
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Can't add block on top of nothing")
	}
}

func (h Handler) getPreviousHash() []byte {
	return Controllers.Hash(h.GetLastBlock())
}

func (h Handler) GetLastBlock() *Models.Block {
	var block Models.Block

	if result := h.DB.Last(&block); result.Error != nil {
		fmt.Println(result.Error)
	}

	return &block
}

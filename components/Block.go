package components

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
)

type Block struct {
	index        int
	nonce        int
	timestamp    int64
	PreviousHash string
	CurrentHash  string
}

func VerifyHash(block Block) bool {
	if hash(block.index, block.nonce, block.timestamp, block.PreviousHash) != block.CurrentHash {
		return false
	}
	return true
}

func MineBlock(index int, nonce int, timestamp int64, previousHash string) Block {
	var t = hash(index, nonce, timestamp, previousHash)
	for t[0:2] != "00" {
		nonce += 1
		t = hash(index, nonce, timestamp, previousHash)
	}
	return Block{index, nonce, timestamp, previousHash, t}
}

func hash(index int, nonce int, timestamp int64, previousHash string) string {
	hasher := sha256.New()
	hasher.Write([]byte(strconv.Itoa(index) + strconv.Itoa(nonce) + strconv.Itoa(int(timestamp)) + previousHash))
	var result string = fmt.Sprintf("%d", hasher.Sum(nil))
	result = strings.ReplaceAll(result, " ", "")
	return result[1:(len(result) - 1)]
}

package database

import (
	"blockchain/components"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func (h handler) GetBlocks(w http.ResponseWriter, r *http.Request) {

}

func (h handler) AddBlock(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var block components.BlockType
	json.Unmarshal(body, &block)

	if result := h.DB.Create(&block); result.Error != nil {
		fmt.Println(result.Error)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Created")
}

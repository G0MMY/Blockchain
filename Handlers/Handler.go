package Handlers

import "blockchain/Models"

type Handler struct {
	Blockchain *Models.Blockchain
}

func New(blockchain *Models.Blockchain) Handler {
	return Handler{blockchain}
}

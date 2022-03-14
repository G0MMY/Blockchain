package main

import (
	"blockchain/Models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func ConnectDatabase() *gorm.DB {
	db, err := gorm.Open(postgres.Open("postgres://postgres:blockchain@localhost:5432/blockchain"), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(&Models.Block{}, &Models.Transaction{}, &Models.Output{}, &Models.Input{}, &Models.MemPoolTransaction{}, &Models.MemPoolInput{}, &Models.MemPoolOutput{}, &Models.Wallet{})

	return db
}

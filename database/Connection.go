package database

import (
	"blockchain/components"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func ConnectDatabase() *gorm.DB {
	db, err := gorm.Open(postgres.Open("postgres://postgres:blockchain@localhost:5432/blockchain"), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(&components.BlockType{})

	return db
}

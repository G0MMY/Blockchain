package Controllers

import (
	"blockchain/Models"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func CreateWallet() Models.Wallet {
	privateKey, publicKey := GenerateNewKeyPair()

	return Models.Wallet{publicKey, privateKey}
}

func GenerateNewKeyPair() ([]byte, []byte) {
	bitSize := 4096

	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		panic(err)
	}

	publicKey := GetPublicKeyFromPrivateKey(privateKey)

	encodedPrivateKey := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)
	encodedPublicKey := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(publicKey.(*rsa.PublicKey)),
		},
	)

	return encodedPrivateKey, encodedPublicKey
}

func GetPublicKeyFromPrivateKey(privateKey *rsa.PrivateKey) crypto.PublicKey {
	return privateKey.Public()
}

func DecodePrivateKey(privateKey []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(privateKey)
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		fmt.Println(err)
	}

	return priv
}

//func Sign(amount int, privateKey []public) []byte {
//
//}
//
//func SignTransaction(privateKey []byte, transaction Models.MemPoolTransaction) {
//	for _, input := range transaction.Inputs {
//		input.Signature = Sign(input.Output.Amount, privateKey)
//	}
//}

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

	return EncodePrivateKey(privateKey), EncodePublicKey(publicKey)
}

func EncodePublicKey(publicKey crypto.PublicKey) []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(publicKey.(*rsa.PublicKey)),
		},
	)
}

func EncodePrivateKey(privateKey *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)
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

func DecryptSignature(publicKey *rsa.PublicKey, signature []byte, amount []byte) error {
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, amount, signature)
}

func Sign(amount []byte, privateKey *rsa.PrivateKey) []byte {
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, amount)
	if err != nil {
		fmt.Println(err)
	}

	return signature
}

//func SignTransaction(privateKey []byte, transaction Models.MemPoolTransaction) {
//	for _, input := range transaction.Inputs {
//		input.Signature = Sign(input.Output.Amount, DecodePrivateKey(privateKey))
//	}
//}

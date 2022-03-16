package Controllers

import (
	"blockchain/Models"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func CreateWallet() Models.Wallet {
	privateKey, publicKey := GenerateNewKeyPair()

	return Models.Wallet{publicKey, privateKey}
}

func GenerateNewKeyPair() (string, string) {
	bitSize := 4096

	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		panic(err)
	}

	publicKey := getPublicKeyFromPrivateKey(privateKey)

	return fmt.Sprintf("%s", EncodePrivateKey(privateKey)), fmt.Sprintf("%s", EncodePublicKey(publicKey))
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

func GetPublicKeyFromPrivateKey(privateKey []byte) []byte {
	priv := DecodePrivateKey(privateKey)
	pub := priv.Public()

	return EncodePublicKey(pub)
}

func getPublicKeyFromPrivateKey(privateKey *rsa.PrivateKey) crypto.PublicKey {
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

func DecodePublicKey(publicKey []byte) *rsa.PublicKey {
	block, _ := pem.Decode(publicKey)
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)

	if err != nil {
		fmt.Println(err)
	}

	return pub
}

func ValidateTransaction(transaction Models.Transaction) bool {
	for _, input := range transaction.Inputs {
		err := validateSignature(DecodePublicKey(input.Output.PublicKey), input.Signature, HashInt(input.Output.Amount))
		if err != nil {
			fmt.Println(err)
			return false
		}
	}
	return true
}

func validateSignature(publicKey *rsa.PublicKey, signature []byte, amount []byte) error {
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, amount, signature)
}

func Sign(amount []byte, privateKey *rsa.PrivateKey) []byte {
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, amount)
	if err != nil {
		fmt.Println(err)
	}

	return signature
}

func HashInt(value int) []byte {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%d", value)))

	return hash[:]
}

func SignTransaction(privateKey []byte, transaction Models.MemPoolTransaction) Models.MemPoolTransaction {
	result := transaction
	for i, input := range transaction.Inputs {
		result.Inputs[i].Signature = Sign(HashInt(input.Output.Amount), DecodePrivateKey(privateKey))
	}

	return result
}

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
	"strings"
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

	stringPrivateKey, stringPublicKey := string(EncodePrivateKey(privateKey)), string(EncodePublicKey(publicKey))

	return CleanPrivateKey(stringPrivateKey), CleanPublicKey(stringPublicKey)
}

func CleanPublicKey(key string) string {
	return strings.ReplaceAll(key[31:len(key)-30], "\n", "")
}

func CleanPrivateKey(key string) string {
	return strings.ReplaceAll(key[32:len(key)-31], "\n", "")
}

func StringPrivateKeyToByte(key string) []byte {
	return []byte("-----BEGIN RSA PRIVATE KEY-----\n" + key + "\n-----END RSA PRIVATE KEY-----")
}
func StringPublicKeyToByte(key string) []byte {
	return []byte("-----BEGIN RSA PUBLIC KEY-----\n" + key + "\n-----END RSA PUBLIC KEY-----")
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

func GetDecodedKey(privateKey []byte) ([]byte, []byte) {
	return EncodePrivateKey(DecodePrivateKey(privateKey)), GetPublicKeyFromPrivateKey(privateKey)
}

func GetPublicKeyFromPrivateKey(privateKey []byte) []byte {
	priv := DecodePrivateKey(privateKey)
	if priv != nil {
		pub := priv.Public()

		return EncodePublicKey(pub)
	}

	return nil
}

func getPublicKeyFromPrivateKey(privateKey *rsa.PrivateKey) crypto.PublicKey {
	return privateKey.Public()
}

func DecodePrivateKey(privateKey []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(privateKey)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		fmt.Println("failed to decode PEM block containing private key")
	} else {
		priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)

		if err != nil {
			fmt.Println(err)
		}
		return priv
	}

	return nil
}

func DecodePublicKey(publicKey []byte) *rsa.PublicKey {
	block, _ := pem.Decode(publicKey)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		fmt.Println("failed to decode PEM block containing public key")
	} else {
		pub, err := x509.ParsePKCS1PublicKey(block.Bytes)

		if err != nil {
			fmt.Println(err)
		}

		return pub
	}
	return nil
}

func ValidateTransaction(transaction Models.MemPoolTransaction) bool {
	for _, input := range transaction.Inputs {
		pub := StringPublicKeyToByte(string(input.Output.PublicKey))
		decodedPub := DecodePublicKey(pub)
		err := validateSignature(decodedPub, input.Signature, HashInt(input.Output.Amount))

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
	fmt.Println(privateKey.Size())
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
		priv := DecodePrivateKey(privateKey)
		if priv != nil {
			result.Inputs[i].Signature = Sign(HashInt(input.Output.Amount), priv)
		}
	}

	return result
}

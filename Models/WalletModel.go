package Models

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"log"
)

const (
	checksumLength = 4
	version        = byte(0x00)
)

type Wallet struct {
	PrivateKey []byte
	PublicKey  []byte
}

func CreateWallet() *Wallet {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)

	if err != nil {
		log.Println(err)
		return &Wallet{}
	}

	publicKey := private.Public()

	return &Wallet{EncodePrivateKey(private), EncodePublicKey(publicKey)}
}

func IsValidPrivateKey(privateKey []byte) bool {
	_, err := x509.ParseECPrivateKey(privateKey)
	if err != nil {
		return false
	}

	return true
}

func IsValidPublicKey(publicKey []byte) bool {
	_, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return false
	}

	return true
}

func EncodePublicKey(public crypto.PublicKey) []byte {
	publicKey, err := x509.MarshalPKIXPublicKey(public)
	if err != nil {
		log.Println(err)
	}

	return publicKey
}

func EncodePrivateKey(private *ecdsa.PrivateKey) []byte {
	privateKey, err := x509.MarshalECPrivateKey(private)
	if err != nil {
		log.Println(err)
		return nil
	}

	return privateKey
}

func DecodePublicKey(publicKey []byte) *ecdsa.PublicKey {
	if !IsValidPublicKey(publicKey) {
		log.Println("Invalid Public Key")
		return nil
	}
	key, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		log.Println(err)
		return nil
	}

	return key.(*ecdsa.PublicKey)
}

func DecodePrivateKey(privateKey []byte) *ecdsa.PrivateKey {
	if !IsValidPrivateKey(privateKey) {
		log.Println("Invalid Private Key")
		return nil
	}
	key, err := x509.ParseECPrivateKey(privateKey)
	if err != nil {
		log.Println(err)
		return nil
	}

	return key
}

func getCheckSum(publicKeyHash []byte) []byte {
	firstHash := sha256.Sum256(publicKeyHash)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}

func GetAddress(publicKey []byte) []byte {
	publicKeyHash := GetPublicKeyHash(publicKey)
	versionHash := append([]byte{version}, publicKeyHash...)
	checkSum := getCheckSum(versionHash)
	fullHash := append(versionHash, checkSum...)
	address := base58.Encode(fullHash)

	return []byte(address)
}

func GetPublicKeyHashFromAddress(address []byte) []byte {
	decodedAddress := base58.Decode(string(address))

	if len(decodedAddress) < checksumLength {
		return []byte{}
	}

	return decodedAddress[1 : len(decodedAddress)-checksumLength]
}

func IsValidAddress(address []byte) bool {
	decodedAddress := base58.Decode(string(address))

	if len(decodedAddress) < checksumLength {
		return false
	}

	checkSum := decodedAddress[len(decodedAddress)-checksumLength:]
	addressVersion := decodedAddress[0]
	publicKeyHash := decodedAddress[1 : len(decodedAddress)-checksumLength]
	versionHash := append([]byte{addressVersion}, publicKeyHash...)
	targetCheckSum := getCheckSum(versionHash)

	return bytes.Compare(checkSum, targetCheckSum) == 0
}

func GetPublicKeyFromPrivateKey(privateKey []byte) []byte {
	if !IsValidPrivateKey(privateKey) {
		log.Println("Invalid private key")
		return nil
	}

	priv := DecodePrivateKey(privateKey)

	return EncodePublicKey(priv.Public())
}

func GetPublicKeyHash(publicKey []byte) []byte {
	shaHash := sha256.Sum256(publicKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(shaHash[:])
	if err != nil {
		log.Println(err)
		return nil
	}

	return hasher.Sum(nil)
}

func ValidateAddress(address []byte) []byte {
	if IsValidPublicKey(address) {
		return GetPublicKeyHash(address)
	} else if IsValidAddress(address) {
		return GetPublicKeyHashFromAddress(address)
	} else {
		log.Println("Invalid address provided")
	}

	return []byte{}
}

func Sign(amount int, privateKey []byte) []byte {
	if !IsValidPrivateKey(privateKey) {
		log.Println("Invalid private key")
		return nil
	}

	priv := DecodePrivateKey(privateKey)

	signature, err := ecdsa.SignASN1(rand.Reader, priv, HashInt(amount))

	if err != nil {
		log.Println(err)
		return nil
	}

	return signature
}

func HashInt(value int) []byte {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%d", value)))

	return hash[:]
}

func ValidateSignature(amount int, publicKey, signature []byte) bool {
	if !IsValidPublicKey(publicKey) {
		log.Println("Invalid public key")
		return false
	}

	return ecdsa.VerifyASN1(DecodePublicKey(publicKey), HashInt(amount), signature)
}

func GetBalance(address []byte, unspentOutputs *UnspentOutput) int {
	balance := 0

	if unspentOutputs != nil && unspentOutputs.Outputs != nil {
		pubKeyHash := ValidateAddress(address)

		for _, output := range unspentOutputs.Outputs {
			if bytes.Compare(pubKeyHash, output.PublicKeyHash) == 0 {
				balance += output.Amount
			}
		}
	}

	return balance
}

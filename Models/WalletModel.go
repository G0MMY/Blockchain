package Models

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"github.com/btcsuite/btcd/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"log"
)

//update PersistUnspentOutputs

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
		log.Panic(err)
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
		log.Panic(err)
	}

	return publicKey
}

func EncodePrivateKey(private *ecdsa.PrivateKey) []byte {
	privateKey, err := x509.MarshalECPrivateKey(private)
	if err != nil {
		log.Panic(err)
	}

	return privateKey
}

func (wallet *Wallet) DecodePublicKey() crypto.PublicKey {
	publicKey, err := x509.ParsePKIXPublicKey(wallet.PublicKey)
	if err != nil {
		log.Panic(err)
	}

	return publicKey.(crypto.PublicKey)
}

func (wallet *Wallet) DecodePrivateKey() *ecdsa.PrivateKey {
	privateKey, err := x509.ParseECPrivateKey(wallet.PrivateKey)
	if err != nil {
		log.Panic(err)
	}

	return privateKey
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

func GetPublicKeyHash(publicKey []byte) []byte {
	shaHash := sha256.Sum256(publicKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(shaHash[:])
	if err != nil {
		log.Panic(err)
	}

	return hasher.Sum(nil)
}

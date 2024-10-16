package main

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"os"
)

func readAddresses(filePath string) (map[string]bool, error) {
	addresses := make(map[string]bool)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		addresses[scanner.Text()] = true
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}

func generateKeyAndAddress() (string, string, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
	}

	publicKey := privateKey.PublicKey
	address, err := publicKeyToAddress(publicKey)
	if err != nil {
		return "", "", err
	}

	return hex.EncodeToString(privateKey.D.Bytes()), address, nil
}

func publicKeyToAddress(publicKey ecdsa.PublicKey) (string, error) {
	pubKeyBytes := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)

	sha256Hash := sha256.New()
	sha256Hash.Write(pubKeyBytes)
	sha256Result := sha256Hash.Sum(nil)

	ripemd160Hash := ripemd160.New()
	ripemd160Hash.Write(sha256Result)
	ripemd160Result := ripemd160Hash.Sum(nil)

	networkVersion := byte(0x00)
	addressBytes := append([]byte{networkVersion}, ripemd160Result...)
	checksum := sha256Checksum(addressBytes)
	fullAddress := append(addressBytes, checksum...)

	return base58.Encode(fullAddress), nil
}

func sha256Checksum(input []byte) []byte {
	firstSHA := sha256.New()
	firstSHA.Write(input)
	result := firstSHA.Sum(nil)

	secondSHA := sha256.New()
	secondSHA.Write(result)
	finalResult := secondSHA.Sum(nil)

	return finalResult[:4]
}

package helper

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
)

func PublicKeyToAddress(publicKey *ecdsa.PublicKey) (address string) {
	return crypto.PubkeyToAddress(*publicKey).Hex()
}

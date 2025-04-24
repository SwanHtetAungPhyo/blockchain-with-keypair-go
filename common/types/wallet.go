package types

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	PubAddress common.Address
}

func NewWallet() *Wallet {
	keys, err := crypto.GenerateKey()
	if err != nil {
		panic(fmt.Sprintf("Failed to generate keys: %v", err))
	}

	address := crypto.PubkeyToAddress(keys.PublicKey)
	//fmt.Printf("Generated wallet address: %s\n", address.Hex())

	return &Wallet{
		PrivateKey: keys,
		PublicKey:  &keys.PublicKey,
		PubAddress: address,
	}
}
func (w *Wallet) SignTx(tx *Transaction) string {
	tx.Id = hex.EncodeToString(crypto.Keccak256([]byte(tx.messageToSign())))

	msg := crypto.Keccak256([]byte(tx.messageToSign()))

	prefixedMsg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(msg), msg)
	prefixedHash := crypto.Keccak256([]byte(prefixedMsg))

	signature, err := crypto.Sign(prefixedHash, w.PrivateKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to sign transaction: %v", err))
	}

	// The signature is now in [R || S || V] format where V is 0 or 1
	// For Ethereum, we need V to be 27 or 28
	signature[64] += 27

	tx.Signature = hex.EncodeToString(signature)
	return tx.Signature
}
func (w *Wallet) Transfer(url string, to common.Address) {

}

func (w *Wallet) Address() string {
	return w.PubAddress.Hex()

}

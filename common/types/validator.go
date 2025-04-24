package types

import (
	"encoding/hex"
	"fmt"
	"github.com/SwanHtetAungPhyo/chain-common/helper"
	"github.com/ethereum/go-ethereum/crypto"
	"sync"
)

type Validator struct {
	Wallet *Wallet // Changed from wallet to Wallet (exported)
	Stake  int     `json:"stake"`
}

func NewValidator() *Validator {
	return &Validator{
		Wallet: NewWallet(), // Properly initialize wallet
		Stake:  0,
	}
}

func (v *Validator) SignBlock(block *Block) {
	if v.Wallet == nil || v.Wallet.PrivateKey == nil {
		panic("validator wallet not initialized")
	}

	hashBytes, err := hex.DecodeString(block.BlockHash)
	if err != nil {
		panic(fmt.Sprintf("invalid block hash hex: %v", err))
	}

	// Create Ethereum signed message
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(hashBytes), hashBytes)
	msgHash := crypto.Keccak256([]byte(msg))

	signature, err := crypto.Sign(msgHash, v.Wallet.PrivateKey)
	if err != nil {
		panic(fmt.Sprintf("failed to sign block: %v", err))
	}

	// Adjust V for Ethereum (27/28)
	signature[64] += 27

	block.ValidatorsSignature = hex.EncodeToString(signature)
}

type CachedValidator struct {
	*Validator
	cachedPubKey []byte
}

var validatorCache sync.Once
var instance *CachedValidator

func NewCachedValidator() *CachedValidator {
	validatorCache.Do(func() {
		v := NewValidator()
		pubKey := crypto.FromECDSAPub(v.Wallet.PublicKey)
		instance = &CachedValidator{
			Validator:    v,
			cachedPubKey: pubKey,
		}
	})
	return instance
}

func (cv *CachedValidator) SignBlock(block *Block) {
	hashBytes := helper.Must(hex.DecodeString(block.BlockHash))

	signature := helper.Must(crypto.Sign(hashBytes, cv.Wallet.PrivateKey))
	
	if signature[64] < 27 {
		signature[64] += 27
	}

	block.ValidatorsSignature = hex.EncodeToString(signature)
}

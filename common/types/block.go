package types

import (
	"encoding/hex"
	"github.com/google/uuid"
	"time"

	"fmt"
	"github.com/SwanHtetAungPhyo/chain-common/helper"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/goccy/go-json"
)

type (
	Header struct {
		MerkleTree [][][]string
		ChainId    string
	}
	Block struct {
		Header              Header         `json:"header"`
		Id                  string         `json:"id"`
		Timestamp           string         `json:"timestamp"`
		MerkleRoot          string         `json:"merkle_root"`
		Txs                 []*Transaction `json:"txs"`
		Validator           *Validator     `json:"-"`
		ValidatorAddress    string         `json:"validator_address"`
		PrevBlockHash       string         `json:"prev_block_hash"`
		BlockHash           string         `json:"block_hash"`
		ValidatorsSignature string         `json:"validators_signature"`
	}
)

func NewBlock(validator *Validator) *Block {
	return &Block{
		Id: uuid.New().String(),
		Header: Header{
			MerkleTree: make([][][]string, 0),
			ChainId:    "0",
		},
		MerkleRoot:       "0",
		Timestamp:        time.Now().Format(time.RFC3339),
		PrevBlockHash:    "",
		Validator:        validator,
		ValidatorAddress: validator.Wallet.Address(),
		Txs:              make([]*Transaction, 0),
	}
}

func (b *Block) AddTx(tx *Transaction) {
	b.Txs = append(b.Txs, tx)
}

func (b *Block) headerCreation() {
	tree, root := helper.BuildMerkleTree(b.allTxToByes())
	for _, level := range tree {
		var levels []string
		for _, hashedNode := range level {
			levels = append(levels, hex.EncodeToString(hashedNode))
		}
		b.Header.MerkleTree = append(b.Header.MerkleTree, [][]string{levels})
	}
	b.MerkleRoot = hex.EncodeToString(root)
}

//func (b *Block) computeBlockHash() {
//	header := fmt.Sprintf("%s|%s|%s|%d",
//		b.PrevBlockHash,
//		b.MerkleRoot,
//		b.Timestamp,
//	)
//	hash := crypto.Keccak256([]byte(header))
//	b.BlockHash = hex.EncodeToString(hash)
//}

func (b *Block) FinalizeBlock() {
	b.headerCreation()
	b.computeBlockHash()

	b.signBlock(b.Validator)

}

func (b *Block) computeBlockHash() {

	headerData := fmt.Sprintf("%s|%s|%s|%s|%d",
		b.PrevBlockHash,
		b.MerkleRoot,
		b.Timestamp,
		b.Validator.Wallet.Address(),
		len(b.Txs),
	)
	b.BlockHash = hex.EncodeToString(crypto.Keccak256([]byte(headerData)))
}

func (b *Block) validate() error {
	if b.BlockHash == "" {
		return fmt.Errorf("missing block hash")
	}
	if b.MerkleRoot == "0" {
		return fmt.Errorf("invalid merkle root")
	}
	if len(b.Txs) == 0 {
		return fmt.Errorf("block contains no transactions")
	}
	if !b.VerifySignature() {
		return fmt.Errorf("invalid validator signature")
	}
	return nil
}

func (b *Block) VerifySignature() bool {
	hashBytes := helper.Must(hex.DecodeString(b.BlockHash))
	sigBytes := helper.Must(hex.DecodeString(b.ValidatorsSignature))

	return crypto.VerifySignature(
		crypto.FromECDSAPub(b.Validator.Wallet.PublicKey),
		hashBytes,
		sigBytes[:64],
	)
}
func (b *Block) constructBlockMessage() []byte {
	return helper.Must(json.Marshal(b))
}

func (b *Block) signBlock(validator *Validator) {
	validator.SignBlock(b)
}

func (b *Block) allTxToByes() [][]byte {
	txBytes := make([][]byte, len(b.Txs))
	for i, tx := range b.Txs {
		txBytes[i] = tx.toBytes()
	}
	return txBytes
}

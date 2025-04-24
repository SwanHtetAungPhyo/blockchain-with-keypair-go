package types

import (
	"fmt"
	"github.com/SwanHtetAungPhyo/chain-common/helper"
	"github.com/goccy/go-json"
	"strings"
	"time"

	eth "github.com/ethereum/go-ethereum/crypto"
)

type (
	Transaction struct {
		Id            string         `json:"id"`
		From          string         `json:"from"`
		To            string         `json:"to"`
		Data          map[string]any `json:"data"`
		Timestamp     string         `json:"timestamp"`
		PrevBlockHash string         `json:"block_hash"`
		BlockId       string         `json:"block_id"`
		Signature     string         `json:"signature"`
	}
)

func NewTransaction(prevBlockHash string, data map[string]any, to string, blockId string, from string) *Transaction {
	return &Transaction{
		PrevBlockHash: prevBlockHash,
		Data:          data,
		To:            to,
		From:          from,
		BlockId:       blockId,
		Timestamp:     time.Now().Format("2006-01-02 15:04:05"),
	}
}

func (tx *Transaction) ToByte() []byte {
	return tx.toBytes()
}

func (tx *Transaction) HashMessage() string {
	messageToSigned := tx.messageToSign()
	hashDataTx := eth.Keccak256Hash([]byte(messageToSigned))
	return hashDataTx.Hex()
}

func (tx *Transaction) mapToString(in map[string]any) string {
	var out strings.Builder
	for k, v := range in {
		out.WriteString(fmt.Sprintf("{%s:%v}", k, v))
	}
	return out.String()
}
func (tx *Transaction) messageToSign() string {
	// Use JSON marshaling with sorted keys
	dataBytes, err := json.Marshal(tx.Data)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal transaction data: %v", err))
	}

	parts := []string{
		tx.Id,
		strings.ToLower(tx.From), // Normalize address case
		strings.ToLower(tx.To),
		string(dataBytes),
		tx.Timestamp,
		tx.PrevBlockHash,
		tx.BlockId,
	}

	return strings.Join(parts, "|")
}
func (tx *Transaction) toBytes() []byte {
	jsonBytes := helper.Must(json.Marshal(tx))
	return jsonBytes
}

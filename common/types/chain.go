package types

import (
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fatih/color"
)

type Chain struct {
	Mu                   sync.RWMutex
	ChainID              string         `json:"chain_id"`
	ChainStarted         string         `json:"chain_started"`
	CompletedBlockNumber int            `json:"completed_block_number"`
	PendingTransactions  []*Transaction `json:"pending_transactions"`
	Blocks               []Block        `json:"blocks"`
	stopChan             chan struct{}
	//txPool               chan *Transaction
}

func NewChain() *Chain {
	chain := &Chain{
		ChainID:              fmt.Sprintf("%x", crypto.Keccak256([]byte("ChainID")))[:8],
		ChainStarted:         time.Now().Format(time.RFC3339),
		CompletedBlockNumber: 0,
		PendingTransactions:  make([]*Transaction, 0),
		Blocks:               make([]Block, 0),
		stopChan:             make(chan struct{}),
		//txPool:               make(chan *Transaction, 1000),
	}
	chain.addGenesisBlock()
	return chain
}
func (c *Chain) addGenesisBlock() {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	genesisBlock := genesisBlockProducer()
	c.Blocks = append(c.Blocks, *genesisBlock)
	c.CompletedBlockNumber = len(c.Blocks)

	color.Green("✓ Genesis block created: %s", genesisBlock.BlockHash[:8]+"...")
}

func (c *Chain) GetLatestHash() string {
	if len(c.Blocks) == 0 {
		return ""
	}
	return c.Blocks[len(c.Blocks)-1].BlockHash
}

func (c *Chain) GetLatestBlockID() *string {
	if len(c.Blocks) == 0 {
		return nil
	}
	return &c.Blocks[len(c.Blocks)-1].Id
}
func (c *Chain) AddBlock(block *Block) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	block.PrevBlockHash = c.Blocks[len(c.Blocks)-1].BlockHash
	c.Blocks = append(c.Blocks, *block)
	c.CompletedBlockNumber = len(c.Blocks)
}

func (c *Chain) GetBlock(height int) *Block {
	if height < 0 || height >= len(c.Blocks) {
		return nil
	}
	return &c.Blocks[height]
}
func (c *Chain) GetBlockByHash(hash string) *Block {
	for _, block := range c.Blocks {
		if block.BlockHash == hash {
			return &block
		}
	}
	return nil
}

//	func (c *Chain) SubmitTransaction(tx *Transaction) {
//		c.txPool <- tx
//	}
//
//	func (c *Chain) SubmitTransaction(tx *Transaction) {
//		c.txPool <- tx
//	}
//
//	func (c *Chain) Start(miningInterval time.Duration) {
//		color.Cyan("\nStarting blockchain with mining interval: %v", miningInterval)
//
//		go c.processTransactions(miningInterval)
//		go c.startStatePrinter()
//
//		color.Green("Blockchain started successfully!")
//	}
//
//	func (c *Chain) Stop() {
//		close(c.stopChan)
//
//		time.Sleep(500 * time.Millisecond)
//
//		c.mu.Lock()
//		defer c.mu.Unlock()
//
//		color.Green("\nFinal Chain State:")
//		color.Cyan("Total Blocks: %d", len(c.Blocks))
//		color.Cyan("Pending Transactions: %d", len(c.PendingTransactions))
//	}
//
//	func (c *Chain) processTransactions(miningInterval time.Duration) {
//		ticker := time.NewTicker(miningInterval)
//		defer ticker.Stop()
//
//		for {
//			select {
//			case tx := <-c.txPool:
//				c.mu.Lock()
//				c.PendingTransactions = append(c.PendingTransactions, tx)
//				pendingCount := len(c.PendingTransactions)
//				c.mu.Unlock()
//
//				color.Blue("↑ New transaction received (Total pending: %d): %s",
//					pendingCount, tx.Id[:8]+"...")
//
//				// Non-blocking check for mining threshold
//				if pendingCount >= 10 {
//					select {
//					case c.txPool <- nil: // Signal to mine
//					default:
//						go c.mineBlock()
//					}
//				}
//
//			case <-ticker.C:
//				c.mu.Lock()
//				pendingCount := len(c.PendingTransactions)
//				c.mu.Unlock()
//
//				if pendingCount > 0 {
//					color.Yellow("⏰ Mining interval reached with %d pending transactions", pendingCount)
//					go c.mineBlock()
//				}
//
//			case <-c.stopChan:
//				color.Yellow("Mining process stopped")
//				return
//			}
//		}
//	}
//
//	func (c *Chain) mineBlock() {
//		c.mu.Lock()
//		defer c.mu.Unlock()
//
//		if len(c.PendingTransactions) == 0 {
//			color.Yellow("No transactions to mine")
//			return
//		}
//
//		txCount := 10
//		if len(c.PendingTransactions) < 10 {
//			txCount = len(c.PendingTransactions)
//		}
//
//		transactions := make([]*Transaction, txCount)
//		copy(transactions, c.PendingTransactions[:txCount])
//		c.PendingTransactions = c.PendingTransactions[txCount:]
//
//		color.Magenta("\n⛏️  Mining new block with %d transactions...", txCount)
//
//		validator := NewValidator()
//		newBlock := NewBlock(validator)
//		if len(c.Blocks) > 0 {
//			newBlock.PrevBlockHash = c.Blocks[len(c.Blocks)-1].BlockHash
//		}
//
//		for _, tx := range transactions {
//			newBlock.AddTx(tx)
//		}
//
//		newBlock.finalizeBlock()
//		c.AddBlock(newBlock)
//
//		color.Green("✅ Block %d mined: %s (Transactions: %d)",
//			c.CompletedBlockNumber,
//			newBlock.BlockHash[:8]+"...",
//			len(newBlock.Txs))
//
//		go c.printState()
//	}
//
//	func (c *Chain) startStatePrinter() {
//		ticker := time.NewTicker(12 * time.Second)
//		defer ticker.Stop()
//
//		for {
//			select {
//			case <-ticker.C:
//				c.printState()
//			case <-c.stopChan:
//				return
//			}
//		}
//	}
//
//	func (c *Chain) printState() {
//		c.mu.Lock()
//		defer c.mu.Unlock()
//
//		f := colorjson.NewFormatter()
//		f.Indent = 4
//
//		state := map[string]interface{}{
//			"chain_id":       c.ChainID,
//			"height":         c.CompletedBlockNumber,
//			"started":        c.ChainStarted,
//			"blocks":         len(c.Blocks),
//			"pending_tx":     len(c.PendingTransactions),
//			"last_block":     c.getLastBlockInfo(),
//			"pending_tx_ids": c.getPendingTxIds(),
//		}
//
//		s, _ := f.Marshal(state)
//		color.Cyan("\n══════════ Blockchain State ══════════")
//		fmt.Println(string(s))
//		color.Cyan("══════════════════════════════════════\n")
//	}
func (c *Chain) getLastBlockInfo() map[string]interface{} {
	if len(c.Blocks) == 0 {
		return nil
	}
	last := c.Blocks[len(c.Blocks)-1]
	return map[string]interface{}{
		"hash":      last.BlockHash[:8] + "...",
		"tx_count":  len(last.Txs),
		"timestamp": last.Timestamp,
		"miner":     last.Validator.Wallet.PubAddress.Hex()[:8] + "...",
	}
}

func (c *Chain) getPendingTxIds() []string {
	ids := make([]string, 0, len(c.PendingTransactions))
	for _, tx := range c.PendingTransactions {
		ids = append(ids, tx.Id[:8]+"...")
	}
	return ids
}

func genesisBlockProducer() *Block {
	genesisWallet := NewWallet()
	genesisTx := gxTxProducer(genesisWallet)

	fmt.Printf("Genesis Tx From: %s\n", genesisTx.From)
	fmt.Printf("Genesis Tx Sig: %s\n", genesisTx.Signature)

	if !VerifyTx(genesisTx) {
		panic("Genesis transaction verification failed")
	}

	validator := NewValidator()
	genesisBlock := NewBlock(validator)
	genesisBlock.AddTx(genesisTx)

	fmt.Println("Finalizing genesis block...")
	genesisBlock.FinalizeBlock()

	return genesisBlock
}
func VerifyTx(gxTx *Transaction) bool {
	sigBytes, err := hex.DecodeString(gxTx.Signature)
	if err != nil || len(sigBytes) != 65 {
		fmt.Printf("Signature decode error: %v\n", err)
		return false
	}

	msgHash := crypto.Keccak256([]byte(gxTx.messageToSign()))

	prefixedMsg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(msgHash), msgHash)
	prefixedHash := crypto.Keccak256([]byte(prefixedMsg))

	// The signature should have V = 27 or 28
	if sigBytes[64] != 27 && sigBytes[64] != 28 {
		fmt.Printf("Invalid V value: %d (should be 27 or 28)\n", sigBytes[64])
		return false
	}

	// For SigToPub, we need V = 0 or 1
	sigBytes[64] -= 27

	pubKey, err := crypto.SigToPub(prefixedHash, sigBytes)
	if err != nil {
		fmt.Printf("SigToPub error: %v\n", err)
		return false
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	fromAddr := common.HexToAddress(gxTx.From)

	return recoveredAddr == fromAddr
}
func gxTxProducer(gxWallet *Wallet) *Transaction {
	gxMap := map[string]any{
		"meta": "hello_from_chain",
	}

	fmt.Printf("Creating genesis tx from address: %s\n", gxWallet.PubAddress.Hex())

	tx := NewTransaction(
		"",
		gxMap,
		gxWallet.PubAddress.Hex(),
		uuid.New().String(),
		gxWallet.PubAddress.Hex(),
	)

	gxWallet.SignTx(tx)
	return tx
}

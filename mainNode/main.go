package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/SwanHtetAungPhyo/chain-common/helper"
	"github.com/TylerBrock/colorjson"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fatih/color"
	"github.com/goccy/go-json"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/SwanHtetAungPhyo/chain-common/types"
)

var printChan = make(chan struct{})
var printWG sync.WaitGroup
var blockSizeSum = 0

func main() {
	source := rand.NewSource(time.Now().UnixNano())
	rand.New(source)
	chain := types.NewChain()
	var txPool []*types.Transaction

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sender := types.NewWallet()
	receiver := types.NewWallet()
	fmt.Printf("Receiver Address: %s\n", receiver.Address())

	for i := 0; i <= 100000; i++ {
		tx := createTransaction(sender.PrivateKey, receiver.Address(), i, chain)
		if tx == nil {
			log.Fatalf("Transaction creation failed at index %d", i)
		}
		fmt.Printf("Submitted transaction %d\n", i+1)
		txPool = append(txPool, tx)
		chain.PendingTransactions = append(chain.PendingTransactions, tx)

	}

	preallocateBlock := make([]*types.Block, 1000000)
	for i := 0; i < 100000; i++ {
		validator := types.NewValidator()
		block := types.NewBlock(validator)
		if block == nil {
			log.Fatalf("Block creation failed at index %d", i)
		}
		preallocateBlock[i] = block
	}
	var totalTransactions int
	var startTime = time.Now()

	for i := 0; i <= len(txPool)-5; i += 5 {
		block := preallocateBlock[i]
		fixedTx := txPool[i : i+5]
		for _, tx := range fixedTx {
			if tx == nil {
				log.Fatalf("Transaction creation failed at index %d", i)
			} else {
				if !types.VerifyTx(tx) {
					log.Println("Transaction verification failed at index %d", i)
				}
				block.AddTx(tx)
				totalTransactions++
			}

		}
		block.FinalizeBlock()
		//for i := range block.Txs {
		//	leaf := helper.HashLeaf(block.Txs[i].ToByte())
		//	merkels := bytes(block.Header.MerkleTree)
		//	proof := helper.GenerateMerkleProof(merkels, i)
		//	if !verifyMerkleProof(leaf, proof, calculateMerkleRoot(merkels)) {
		//		fmt.Println("Merkle Proof Verification failed! Cannot produce block.")
		//		return
		//	}
		//}
		chain.AddBlock(block)
		jsonBlock := helper.Must(json.Marshal(block))
		blockSizeSum += len(jsonBlock)
		chain.CompletedBlockNumber = len(chain.Blocks)
		chain.PendingTransactions = nil

	}

	elapsedTime := time.Since(startTime)
	dataThroughput := float64(blockSizeSum) / elapsedTime.Seconds()

	transactionThroughput := float64(totalTransactions) / elapsedTime.Seconds()

	printChainStatus(chain)
	jsondata := helper.Must(json.MarshalIndent(chain, "", "  "))
	file := helper.Must(os.Create("chain.json"))
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err.Error())
		}
	}(file)
	_, err := file.Write(jsondata)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Transaction Throughput: %.2f transactions per second\n", transactionThroughput)
	fmt.Printf("Data Throughput: %.2f bytes per second\n", dataThroughput)
	//printWG.Add(1)
	//go func() {
	//	defer printWG.Done()
	//	ticker := time.NewTicker(10 * time.Second)
	//	defer ticker.Stop()
	//
	//	for {
	//		select {
	//		case <-ticker.C:
	//			printWholeChain(chain)
	//		case <-sigChan:
	//			return
	//		}
	//	}
	//}()
	//// Mining process
	//validator := types.NewValidator()
	//go func() {
	//	for {
	//		select {
	//		case <-sigChan:
	//			return
	//		case <-time.After(10 * time.Second):
	//			printChainStatus(chain)
	//		default:
	//			mineBlock(chain, validator, txPool)
	//			time.Sleep(5 * time.Second)
	//		}
	//	}
	//}()
	//
	//<-sigChan
	//fmt.Println("\nShutting down...")
	////printChainStatus(chain)
	//printWholeChain(chain)
	//fmt.Println("Clean shutdown completed")
}
func verifyMerkleProof(leaf []byte, proof [][]byte, root []byte) bool {
	currentHash := leaf
	for _, siblingHash := range proof {
		currentHash = hashPair(currentHash, siblingHash)
	}
	return string(currentHash) == string(root)
}
func hashPair(left, right []byte) []byte {

	combined := append(left, right...)
	return crypto.Keccak256(combined)
}
func calculateMerkleRoot(tree [][][]byte) []byte {

	currentLevel := tree[len(tree)-1]

	for len(currentLevel) > 1 {
		var parentLevel [][]byte
		for i := 0; i < len(currentLevel); i += 2 {
			var hash []byte
			if i+1 < len(currentLevel) {
				hash = append(currentLevel[i], currentLevel[i+1]...)
			} else {
				hash = append(currentLevel[i], currentLevel[i]...)
			}
			parentLevel = append(parentLevel, hash)
		}
		currentLevel = parentLevel
	}

	return currentLevel[0]
}
func bytes(strs [][][]string) [][][]byte {
	var byteTree [][][]byte
	for _, level := range strs {
		var byteLevel [][]byte
		for _, item := range level {
			var byteItem []byte
			for _, str := range item {
				byteItem = append(byteItem, []byte(str)...)
			}
			byteLevel = append(byteLevel, byteItem)
		}
		byteTree = append(byteTree, byteLevel)
	}
	return byteTree
}

func createTransaction(senderKey *ecdsa.PrivateKey, receiver string, nonce int, chain *types.Chain) *types.Transaction {
	tx := types.NewTransaction(
		chain.GetLatestHash(),
		map[string]interface{}{
			"type":  "transfer",
			"from":  crypto.PubkeyToAddress(senderKey.PublicKey).Hex(),
			"to":    receiver,
			"value": rand.Intn(1000) + 1,
		},
		chain.ChainID,
		receiver,
		crypto.PubkeyToAddress(senderKey.PublicKey).Hex(),
	)

	if tx == nil {
		log.Println("tx is nil after NewTransaction")
		return nil
	}

	sender := &types.Wallet{PrivateKey: senderKey}

	log.Println("Attempting to sign transaction")
	sender.SignTx(tx)

	return tx
}

//	func mineBlock(chain *types.Chain, validator *types.Validator, txPool <-chan *types.Transaction) {
//		block := types.NewBlock(validator)
//		txCount := 0
//
//		for i := 0; i < 5; i++ {
//			select {
//			case tx := <-txPool:
//				if types.VerifyTx(tx) {
//					block.AddTx(tx)
//					txCount++
//				}
//			default:
//				break
//			}
//		}
//
//		if txCount > 0 {
//			block.FinalizeBlock()
//			chain.AddBlock(block)
//			fmt.Printf("\n✅ Mined new block (%s) with %d transactions\n",
//				block.BlockHash[:8],
//				txCount)
//		}
//	}
func printChainStatus(chain *types.Chain) {
	fmt.Printf("\n📦 Final Chain Status:")
	fmt.Printf("\n- Total Blocks: %d", len(chain.Blocks))
	fmt.Printf("\n- Pending Transactions: %d", len(chain.PendingTransactions))
	fmt.Printf("\n- Completed Blocks: %d", chain.CompletedBlockNumber)

	chainSizeBytes, err := json.Marshal(chain)
	if err != nil {
		log.Fatal(err)
	}

	chainSizeMB := float64(len(chainSizeBytes)) / (1024 * 1024)
	avgBlockSizeMB := 0.0
	if len(chain.Blocks) > 0 {
		avgBlockSizeMB = float64(blockSizeSum) / float64(len(chain.Blocks)) / (1024 * 1024)
	}

	fmt.Printf("\n- Chain Size: %.2f MB", chainSizeMB)
	fmt.Printf("\n- Average Block Size: %.4f MB", avgBlockSizeMB)

	if len(chain.Blocks) > 0 {
		lastBlock := chain.Blocks[len(chain.Blocks)-1]
		fmt.Printf("\n- Last Block: %s (%d transactions)",
			lastBlock.BlockHash[:8],
			len(lastBlock.Txs))
	}
	fmt.Println()
}

func printWholeChain(c *types.Chain) {
	c.Mu.RLock()
	defer c.Mu.RUnlock()

	f := colorjson.NewFormatter()
	f.Indent = 4
	f.KeyColor = color.New(color.FgHiCyan)
	f.StringColor = color.New(color.FgHiGreen)
	f.NumberColor = color.New(color.FgHiMagenta)

	s, _ := f.Marshal(c.Blocks)
	color.Cyan("\n══════════ Blockchain State ══════════")
	fmt.Println(string(s))
	color.Cyan("══════════════════════════════════════\n")
}

func formatBlocks(blocks []types.Block) []map[string]interface{} {
	formatted := make([]map[string]interface{}, len(blocks))
	for i, block := range blocks {
		formatted[i] = map[string]interface{}{
			"height":      i,
			"hash":        block.BlockHash[:8] + "...",
			"timestamp":   block.Timestamp,
			"tx_count":    len(block.Txs),
			"prev_hash":   block.PrevBlockHash + "...",
			"tx":          block.Txs,
			"merkle_root": block.MerkleRoot + "...",
			"validator":   block.Validator.Wallet.Address()[:8] + "...",
		}
	}
	return formatted
}

func getLastBlockInfo(c *types.Chain) interface{} {
	if len(c.Blocks) == 0 {
		return "none"
	}
	last := c.Blocks[len(c.Blocks)-1]
	return map[string]interface{}{
		"hash":      last.BlockHash[:8] + "...",
		"tx_count":  len(last.Txs),
		"timestamp": last.Timestamp,
	}
}

/*
      1234
	12  34
1 2     3 4
*/

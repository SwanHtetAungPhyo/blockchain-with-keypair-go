package main

//import (
//	"crypto/ecdsa"
//	"crypto/sha256"
//	"encoding/hex"
//	"encoding/json"
//	"fmt"
//	"math/rand"
//	"strconv"
//	"time"
//
//	"github.com/ethereum/go-ethereum/crypto"
//	"github.com/google/uuid"
//)
//
//const (
//	Pending = "pending"
//	Failed  = "failed"
//	Success = "success"
//)
//
//type Block struct {
//	Index        int
//	Timestamp    string
//	MerkleRoot   []byte
//	Transactions [][]byte
//}
//type Transaction struct {
//	Id        string `json:"id"`
//	From      string `json:"from"`
//	To        string `json:"to"`
//	Timestamp string `json:"timestamp"`
//	Sign      string `json:"sign"`
//	PubKey    string `json:"pub_key"` // hex-encoded public key
//	Status    string `json:"status"`
//}
//
//func must[T any](v T, err error) T {
//	if err != nil {
//		panic(err)
//	}
//	return v
//}
//
//func (t Transaction) toByte() []byte {
//	return must(json.Marshal(&t))
//}
//
//func hashConcat(a, b []byte) []byte {
//	hash := sha256.Sum256(append(a, b...))
//	return hash[:]
//}
//
//func hashLeaf(data []byte) []byte {
//	hash := sha256.Sum256(data)
//	return hash[:]
//}
//
//func buildMerkleTree(transactions [][]byte) ([][][]byte, []byte) {
//	if len(transactions) == 0 {
//		return nil, nil
//	}
//
//	var tree [][][]byte
//	current := make([][]byte, len(transactions))
//
//	for i, tx := range transactions {
//		current[i] = hashLeaf(tx)
//	}
//
//	tree = append(tree, current)
//	fmt.Println("📦 Initial Leaf Hashes:")
//	for i, hash := range current {
//		fmt.Printf("  Leaf %d: %s\n", i, hex.EncodeToString(hash))
//	}
//	fmt.Println()
//
//	levelCount := 0
//	for len(current) > 1 {
//		var nextLevel [][]byte
//		fmt.Printf("🌲 Level %d:\n", levelCount+1)
//
//		for i := 0; i < len(current); i += 2 {
//			left := current[i]
//			right := current[i]
//			if i+1 < len(current) {
//				right = current[i+1]
//			}
//			parent := hashConcat(left, right)
//			nextLevel = append(nextLevel, parent)
//
//			fmt.Printf("  Node[%d]:\n    Left : %s\n    Right: %s\n    -> Parent: %s\n", i/2,
//				hex.EncodeToString(left),
//				hex.EncodeToString(right),
//				hex.EncodeToString(parent))
//		}
//		fmt.Println()
//
//		tree = append(tree, nextLevel)
//		current = nextLevel
//		levelCount++
//	}
//
//	root := tree[len(tree)-1][0]
//	fmt.Println("🌟 Merkle Root:", hex.EncodeToString(root))
//	return tree, root
//}
//
//func generateMerkleProof(tree [][][]byte, index int) [][]byte {
//	var proof [][]byte
//	for level := 0; level < len(tree)-1; level++ {
//		sibling := index ^ 1
//		var siblingHash []byte
//		if sibling < len(tree[level]) {
//			siblingHash = tree[level][sibling]
//		} else {
//			siblingHash = tree[level][index]
//		}
//		proof = append(proof, siblingHash)
//		index /= 2
//	}
//	return proof
//}
//
//func VerifyMerkleProof(leaf []byte, index int, root []byte, proof [][]byte) bool {
//	hash := leaf
//	for _, sibling := range proof {
//		if index%2 == 0 {
//			hash = hashConcat(hash, sibling)
//		} else {
//			hash = hashConcat(sibling, hash)
//		}
//		index = index / 2
//	}
//	return hex.EncodeToString(hash) == hex.EncodeToString(root)
//}
//
//func generateKeypair() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
//	priv := must(crypto.GenerateKey())
//	return priv, &priv.PublicKey
//}
//
//func main() {
//	rand.Seed(time.Now().UnixNano())
//	randCount := rand.Intn(10) + 5
//	transactions := make([]Transaction, randCount)
//
//	for i := 0; i < 5; i++ {
//		priv, pub := generateKeypair()
//		from := crypto.PubkeyToAddress(*pub).Hex()
//
//		tx := Transaction{
//			Id:        uuid.New().String(),
//			From:      from,
//			To:        strconv.Itoa(rand.Intn(1000)),
//			Timestamp: time.Now().UTC().Format(time.RFC3339),
//			Status:    Pending,
//		}
//
//		msg := fmt.Sprintf(`{"from":"%s","to":"%s","timestamp":"%s"}`, tx.From, tx.To, tx.Timestamp)
//		hash := crypto.Keccak256([]byte(msg))
//		sig := must(crypto.Sign(hash, priv))
//
//		tx.Sign = hex.EncodeToString(sig)
//		tx.PubKey = hex.EncodeToString(crypto.FromECDSAPub(pub))
//		transactions[i] = tx
//	}
//
//	var verified [][]byte
//
//	for i, tx := range transactions {
//		msg := fmt.Sprintf(`{"from":"%s","to":"%s","timestamp":"%s"}`, tx.From, tx.To, tx.Timestamp)
//		hash := crypto.Keccak256([]byte(msg))
//
//		pubKeyBytes, err := hex.DecodeString(tx.PubKey)
//		if err != nil {
//			tx.Status = Failed
//			transactions[i] = tx
//			continue
//		}
//
//		sigBytes, err := hex.DecodeString(tx.Sign)
//		if err != nil || len(sigBytes) < 65 {
//			tx.Status = Failed
//			transactions[i] = tx
//			continue
//		}
//
//		valid := crypto.VerifySignature(pubKeyBytes, hash, sigBytes[:64])
//		if valid {
//			tx.Status = Success
//			verified = append(verified, tx.toByte())
//		} else {
//			tx.Status = Failed
//		}
//		transactions[i] = tx
//		fmt.Printf("Tx ID: %s | From: %s | Status: %s\n", tx.Id, tx.From, tx.Status)
//	}
//
//	if len(verified) == 0 {
//		fmt.Println("No valid transactions to build a Merkle tree.")
//		return
//	}
//
//	tree, root := buildMerkleTree(verified)
//	fmt.Println("Merkle Root:", hex.EncodeToString(root))
//
//	block := Block{
//		Index:        0,
//		Timestamp:    time.Now().UTC().Format(time.RFC3339),
//		MerkleRoot:   root,
//		Transactions: verified,
//	}
//
//	for i := range block.Transactions {
//		leaf := hashLeaf(block.Transactions[i]) // Corrected: Use the hash of the transaction data
//		proof := generateMerkleProof(tree, i)
//		if !VerifyMerkleProof(leaf, i, block.MerkleRoot, proof) {
//			fmt.Println("Merkle Proof Verification failed! Cannot produce block.")
//			return
//		}
//	}
//
//	fmt.Println("Merkle Proof Verification Success!")
//}

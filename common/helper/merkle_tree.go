package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func HashConcat(a, b []byte) []byte {
	hash := sha256.Sum256(append(a, b...))
	return hash[:]
}

func HashLeaf(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}
func BuildMerkleTree(transactions [][]byte) ([][][]byte, []byte) {
	if len(transactions) == 0 {
		return nil, nil
	}

	var tree [][][]byte
	current := make([][]byte, len(transactions))

	for i, tx := range transactions {
		current[i] = HashLeaf(tx)
	}

	tree = append(tree, current)
	fmt.Println("📦 Initial Leaf Hashes:")
	for i, hash := range current {
		fmt.Printf("  Leaf %d: %s\n", i, hex.EncodeToString(hash))
	}
	fmt.Println()

	levelCount := 0
	for len(current) > 1 {
		var nextLevel [][]byte
		fmt.Printf("🌲 Level %d:\n", levelCount+1)

		for i := 0; i < len(current); i += 2 {
			left := current[i]
			right := current[i]
			if i+1 < len(current) {
				right = current[i+1]
			}
			parent := HashConcat(left, right)
			nextLevel = append(nextLevel, parent)

			fmt.Printf("  Node[%d]:\n    Left : %s\n    Right: %s\n    -> Parent: %s\n", i/2,
				hex.EncodeToString(left),
				hex.EncodeToString(right),
				hex.EncodeToString(parent))
		}
		fmt.Println()

		tree = append(tree, nextLevel)
		current = nextLevel
		levelCount++
	}

	root := tree[len(tree)-1][0]
	fmt.Println("🌟 Merkle Root:", hex.EncodeToString(root))
	return tree, root
}

func GenerateMerkleProof(tree [][][]byte, index int) [][]byte {
	var proof [][]byte
	for level := 0; level < len(tree)-1; level++ {
		sibling := index ^ 1
		var siblingHash []byte
		if sibling < len(tree[level]) {
			siblingHash = tree[level][sibling]
		} else {
			siblingHash = tree[level][index]
		}
		proof = append(proof, siblingHash)
		index /= 2
	}
	return proof
}

func VerifyMerkleProof(leaf []byte, index int, root []byte, proof [][]byte) bool {
	hash := leaf
	for _, sibling := range proof {
		if index%2 == 0 {
			hash = HashConcat(hash, sibling)
		} else {
			hash = HashConcat(sibling, hash)
		}
		index = index / 2
	}
	return hex.EncodeToString(hash) == hex.EncodeToString(root)
}

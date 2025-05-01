
---

## 🔗 Basic Blockchain with Go

This project is a **simple simulation of a blockchain network** implemented in Go. It demonstrates core blockchain concepts including digital signatures, a mempool, block mining, and Merkle tree-based integrity verification.


I am currently working on the version of this one with Rest API for the client connection . 
GRPC for the blockchain State
Valiator node to verify the transaction and propose the block to the network
### 🚀 Features

- Users **sign messages** using their Ethereum wallet's private key.
- Signed messages are **sent to the network**.
- The network **verifies signatures** and **stores valid transactions** in the **mempool**.
- Once the mempool is full:
    - The network **mines a new block**.
    - Each block includes the **Merkle root** of all transactions.
    - The block is **appended to the blockchain**.
- A **Merkle tree** is used to verify the **integrity** of each block.

### 🧠 Concepts Demonstrated

- Digital signatures with Ethereum private keys
- Transaction queueing (mempool)
- Block mining trigger mechanism
- Merkle tree construction and validation
- Simplified consensus (single miner simulation)

## Run the Project
  ````
  git clone <url>
  cd mainNode
  go run main.go
   ````
  output: the snapshot of the blockchain will be displayed and stored in chain.json

### 📷 Visual Overview

![img_1.png](img_1.png)


---

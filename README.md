# MiniBTC

A lightweight Bitcoin implementation from scratch in Go, inspired by Bitcoin and Decred protocols.

## 🚀 Overview

MiniBTC is a educational cryptocurrency implementation that demonstrates core blockchain concepts including proof-of-work consensus, UTXO model, digital signatures, and peer-to-peer networking. Built entirely in Go without external blockchain libraries, this project showcases fundamental cryptographic and distributed systems principles.

## ✨ Features

- **Blockchain Core**
  - Proof-of-Work consensus mechanism
  - UTXO (Unspent Transaction Output) model
  - Block validation and chain reorganization
  - Merkle tree transaction verification

- **Cryptography**
  - ECDSA digital signatures
  - SHA-256 hashing
  - Address generation and validation
  - Transaction signing and verification

- **Networking**
  - Peer-to-peer protocol implementation
  - Node discovery and connection management
  - Block and transaction propagation
  - Network message handling

- **Transaction System**
  - Input/output transaction model
  - Transaction fees and validation
  - Multi-signature support foundation
  - Mempool management

## 🏗️ Architecture

```
├── blockchain/     # Core blockchain logic
├── crypto/         # Cryptographic utilities
├── network/        # P2P networking layer
├── transaction/    # Transaction handling
├── wallet/         # Wallet functionality
├── consensus/      # Consensus algorithms
└── utils/          # Helper utilities
```

## 🛠️ Installation

```bash
# Clone the repository
git clone https://github.com/Fantasim/minibtc.git
cd minibtc

# Install dependencies
go mod tidy

# Build the project
go build -o minibtc ./cmd/minibtc
```

## 🚀 Quick Start

### Running a Node

```bash
# Start a mining node
./minibtc node --mine

# Start a regular node
./minibtc node

# Connect to specific peers
./minibtc node --peers=127.0.0.1:8333,192.168.1.100:8333
```

### Wallet Operations

```bash
# Create a new wallet
./minibtc wallet create

# Generate a new address
./minibtc wallet address

# Check balance
./minibtc wallet balance

# Send transaction
./minibtc wallet send --to=<address> --amount=<amount>
```

### Mining

```bash
# Start mining
./minibtc mine --address=<mining-address>

# Set mining difficulty
./minibtc mine --difficulty=4
```

## 🔧 Configuration

Create a `config.json` file:

```json
{
  "network": {
    "port": 8333,
    "max_peers": 8,
    "protocol_version": 1
  },
  "mining": {
    "difficulty": 4,
    "block_time": 600
  },
  "wallet": {
    "data_dir": "./wallet_data"
  }
}
```

## 📊 Technical Specifications

| Feature | Implementation |
|---------|----------------|
| **Hashing Algorithm** | SHA-256 |
| **Signature Scheme** | ECDSA (secp256k1) |
| **Consensus** | Proof-of-Work |
| **Block Time** | ~10 minutes (configurable) |
| **Max Block Size** | 1MB |
| **Address Format** | Base58Check encoded |

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./blockchain/...
```

##  Key Concepts 

### Blockchain Fundamentals
- **Block Structure**: Headers, transactions, merkle roots
- **Chain Validation**: Longest chain rule, orphan handling
- **Difficulty Adjustment**: Dynamic mining difficulty

### Cryptographic Security
- **Digital Signatures**: Transaction authorization
- **Hash Functions**: Block linking and proof-of-work
- **Key Management**: Public/private key pairs

### Network Protocol
- **Message Types**: Block, transaction, ping/pong
- **Peer Discovery**: Bootstrap nodes, peer exchange
- **Consensus**: Block propagation and validation

### Transaction Model
- **UTXO Set**: Unspent transaction tracking
- **Script System**: Basic transaction scripting
- **Fee Market**: Transaction prioritization


## 🔍 Code Structure Highlights

```go
// Example: Block structure
type Block struct {
    Header       BlockHeader
    Transactions []Transaction
}

type BlockHeader struct {
    PrevHash     [32]byte
    MerkleRoot   [32]byte
    Timestamp    int64
    Difficulty   uint32
    Nonce        uint32
}
```

## 🚧 Development Status

- ✅ Core blockchain functionality
- ✅ Basic wallet operations
- ✅ Proof-of-work mining
- ✅ P2P networking
- 🔄 Advanced scripting 


## 🙏 Acknowledgments

- **Bitcoin**: Original cryptocurrency design and implementation
- **Decred**: Improved consensus mechanisms and governance
- **Go Community**: Excellent tooling and libraries

package blockchain

type Block struct {
	Size [4]byte
	Header BlockHeader
	Counter uint
	Transactions []Transaction
}

type BlockHeader struct{
	Version [4]byte
	HashPrevBlock [32]byte
	HashMerkleRoot [32]byte
	Time [4]byte
	Bits [4]byte
	Nonce [4]byte
}

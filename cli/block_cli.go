package cli

import (
	"tway/wire"
	"tway/util"
	"tway/script"
	"tway/wallet"
	b "tway/blockchain"
	"fmt"
	"flag"
	"encoding/hex"
	"log"
)

func BlockPrintUsage(){
	fmt.Println(" Options:")
	fmt.Println(" --hash \t block's hash")
	fmt.Println(" --new \t Create and add new blockchain onto the blockchain")
}


func printBlockInChain(block *wire.Block, height int){
	fmt.Printf("Block height: %d\n", height)
	fmt.Printf("Hash: %x\n", block.GetHash())
	fmt.Printf("Merkle root: %x\n", block.Header.HashMerkleRoot)
	fmt.Printf("Size: %d\n", util.DecodeInt(block.Size))
	fmt.Printf("Unix time: %d\n", util.DecodeInt(block.Header.Time))
	fmt.Printf("Difficulty: %d\n", util.DecodeInt(block.Header.Bits))
	fmt.Printf("Txs: %d\n", len(block.Transactions))
	for idx, tx := range block.Transactions {
		fmt.Printf("\t=== Tx [%d] ===\n", idx)
		fmt.Printf("\t Coinbase: %t\n", tx.IsCoinbase())
		fmt.Printf("\t Hash: %x\n", tx.GetHash())
		for idx, out := range tx.Outputs {
			fmt.Printf("\t output[%d] Value: %d\n", idx, util.DecodeInt(out.Value))
			fmt.Printf("\t output[%d] scriptPubKey: %s\n", idx, script.Script.String(out.ScriptPubKey))
		}
	}
}

func printBlock(block *wire.Block){
	fmt.Printf("Hash: %x\n", block.GetHash())
	fmt.Printf("Merkle root: %x\n", block.Header.HashMerkleRoot)
	fmt.Printf("Size: %d\n", util.DecodeInt(block.Size))
	fmt.Printf("Unix time: %d\n", util.DecodeInt(block.Header.Time))
	fmt.Printf("Difficulty: %d\n", util.DecodeInt(block.Header.Bits))
	fmt.Printf("Txs: %d\n", len(block.Transactions))
	for idx, tx := range block.Transactions {
		fmt.Printf("\t=== Tx [%d] ===\n", idx)
		fmt.Printf("\t Coinbase: %t\n", tx.IsCoinbase())
		fmt.Printf("\t Hash: %x\n", tx.GetHash())
		for idx, out := range tx.Outputs {
			fmt.Printf("\t output[%d] Value: %d\n", idx, util.DecodeInt(out.Value))
			fmt.Printf("\t output[%d] scriptPubKey: %s\n", idx, script.Script.String(out.ScriptPubKey))
		}
	}
}

func NewBlock(txs []wire.Transaction, fees int){
	block := wire.NewBlock(txs, b.BC.Tip, wallet.RandomWallet().PublicKey, fees)
	//Créer une target de proof of work
	pow := b.NewProofOfWork(block)
	//cherche le nonce correspondant à la target
	nonce, _, err := pow.Run()
	if err != nil {
		log.Panic(err)
	}
	//ajoute le nonce au header
	block.Header.Nonce = util.EncodeInt(nonce)
	//ajoute la taille total du block
	block.Size = util.EncodeInt(int(block.GetSize()))
	if err := b.BC.AddBlock(block); err != nil {
		fmt.Println("Block non miné")
	}
}

func BlockPrintCli(){
	blockCMD := flag.NewFlagSet("block", flag.ExitOnError)
	hash := blockCMD.String("hash", "", "Print block if exist")
	new := blockCMD.Bool("new", false, "Create and mine new block")
	handleParsingError(blockCMD)

	if *hash != "" {
		h, _ := hex.DecodeString(*hash)
		block, height := b.BC.GetBlockByHash(h)
		if height != -1 {
			printBlockInChain(block, height)
		}
	} else if *new == true {
		var empty []wire.Transaction
		NewBlock(empty, 0)
	} else {
		BlockPrintUsage()
	}
}
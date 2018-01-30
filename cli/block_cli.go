package cli

import (
	"tway/wire"
	"fmt"
	"tway/util"
	"flag"
	"tway/script"
	"encoding/hex"
	b "tway/blockchain"
)

func BlockPrintUsage(){
	fmt.Println(" Options:")
	fmt.Println(" --hash \t block's hash")
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

func BlockPrintCli(){
	blockCMD := flag.NewFlagSet("block", flag.ExitOnError)
	hash := blockCMD.String("hash", "", "Print block if exist")
	handleParsingError(blockCMD)

	if *hash != "" {
		h, _ := hex.DecodeString(*hash)
		block, height := b.BC.GetBlockByHash(h)
		if height != -1 {
			printBlockInChain(block, height)
		}
	} else {
		BlockPrintUsage()
	}
}
package cli

import (
	"tway/twayutil"
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
	fmt.Println(" --hash \t Print a block by its hash")
	fmt.Println(" --height \t Print a block by its height")
	fmt.Println(" --last \t Print last block")
	fmt.Println(" --loop \t Loop execution of a cmd. /!\""+"Works only with : --new and --remove")
	fmt.Println(" --new \t Create and add new blockchain onto the blockchain")
	fmt.Println(" --remove \t Remove block. /!\""+"Works only with --last")
}


func printBlockInChain(block *twayutil.Block, height int){
	fmt.Printf("Block height: %d\n", height)
	fmt.Printf("Hash: %x\n", block.GetHash())
	fmt.Printf("Prev: %x\n", block.Header.HashPrevBlock)
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

func printBlock(block *twayutil.Block){
	fmt.Printf("Hash: %x\n", block.GetHash())
	fmt.Printf("Prev: %x\n", block.Header.HashPrevBlock)
	fmt.Printf("Unix time: %d\n", util.DecodeInt(block.Header.Time))
	fmt.Printf("Difficulty: %d\n", util.DecodeInt(block.Header.Bits))
	fmt.Printf("Txs: %d\n", len(block.Transactions))
	for idx, tx := range block.Transactions {
		fmt.Printf("\t=== Tx [%d] ===\n", idx)
		fmt.Printf("\t Coinbase: %t\n", tx.IsCoinbase())
		fmt.Printf("\t Hash: %x\n", tx.GetHash())
		for idx, out := range tx.Outputs {
			fmt.Printf("\t output[%d] Value: %d\n", idx, util.DecodeInt(out.Value))
		}
	}
}

func NewBlock(txs []twayutil.Transaction, fees int){
	block := twayutil.NewBlock(txs, b.BC.Tip, wallet.NewMiningWallet(), fees, b.BC.GetNewBits())
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


func lastCMD(remove, loop bool){
	if remove == true {
		for {
			_, err := b.BC.RemoveLastBlock()
			if err != nil {
				return
				fmt.Println(err)
			} else {
				fmt.Println("block [",b.BC.Height + 1, "] successfully removed")
			}
			if loop == false {
				return
			}
		}
	} else {
		block := b.BC.GetLastBlock()
		printBlockInChain(block, b.BC.Height)
	}
}

func heightCMD(height int){
	block := b.BC.GetBlockByHeight(height)
	if block != nil {
		printBlock(block)
	} else {
		fmt.Println("This block height doesn't exist.")
	}
}

func newCMD(loop bool){
	var empty []twayutil.Transaction
	for {
		NewBlock(empty, 0)
		if loop == false {
			return
		}
	}
}

func hashCMD(hash string){
	h, _ := hex.DecodeString(hash)
	block, height := b.BC.GetBlockByHash(h)
	if height != -1 {
		printBlockInChain(block, height)
	}
}

func BlockPrintCli(){
	blockCMD := flag.NewFlagSet("block", flag.ExitOnError)
	hash := blockCMD.String("hash", "", "Print block if exist")
	new := blockCMD.Bool("new", false, "Create and mine new block")
	loop := blockCMD.Bool("loop", false, "Loop execution of a cmd. /!\""+"Works only with --new")
	last := blockCMD.Bool("last", false, "print last block")
	height := blockCMD.Int("height", 0, "Print a block by its height.")
	remove := blockCMD.Bool("remove", false, "remove block. /!\""+"Works only with --last")

	handleParsingError(blockCMD)

	if *hash != "" {
		hashCMD(*hash)
	} else if *new == true {
		newCMD(*loop)
	} else if *last == true {
		lastCMD(*remove, *loop)
	} else if *height > 0 {
		heightCMD(*height)
	} else {
		BlockPrintUsage()
	}
}
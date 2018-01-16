package cli

import (
	"fmt"
	"flag"
	"letsgo/blockchain"
	"letsgo/util"
	"encoding/hex"
)

func TxPrintUsage(){
	fmt.Println(" Options:")
	fmt.Println(" --hash \t Print tx equal to hash")
	fmt.Println("Others cmds starting by tx :")
	fmt.Println("\t tx_create")
}

func TxCreateUsage(){
	fmt.Println(" Options:")
	fmt.Println(" --to \t address to send")
	fmt.Println(" --amount \t amount to send")
}

func printTx(tx *blockchain.Transaction, block *blockchain.Block, height int){
	fmt.Printf("Block height: %d\n", height)
	fmt.Printf("Block hash: %x\n\n", block.GetHash())
	fmt.Printf("== TX %x ==\n", tx.GetHash())
	fmt.Printf("    Coinbase: %t\n", tx.IsCoinbase())
	fmt.Printf("    Version: %x\n", tx.Version)
	fmt.Printf("    Value %d\n\n", tx.GetValue())
	fmt.Printf("    %d inputs:\n\n", len(tx.Inputs))
	for idx, in := range tx.Inputs {
		fmt.Printf("    === [%d] ===\n", idx)
		fmt.Printf("    PrevHash: %x\n", in.PrevTransactionHash)
		fmt.Printf("    Vout: %d\n", util.DecodeInt(in.Vout))
	}
	fmt.Printf("    %d outputs:\n\n", len(tx.Outputs))
	for idx, out := range tx.Outputs {
		fmt.Printf("    === [%d] ===\n", idx)
		fmt.Printf("    Value: %d\n", util.DecodeInt(out.Value))
	}
}

func TxPrintCli(){
	TxCMD := flag.NewFlagSet("tx", flag.ExitOnError)
	hash := TxCMD.String("hash", "", "Print tx if exist")
	handleParsingError(TxCMD)

	if *hash != "" {
		h, _ := hex.DecodeString(*hash)
		tx, block, height := blockchain.GetTxByHash(h)
		if height != -1 {
			printTx(tx, block, height)
		}
	} else {
		TxPrintUsage()
	}
}

func TxCreateCli(){
	TxCMD := flag.NewFlagSet("tx_create", flag.ExitOnError)
	to := TxCMD.String("to", "", "address to send")
	amount := TxCMD.Int("amount", 0, "amount to send")
	handleParsingError(TxCMD)

	if *to != "" && *amount > 0 {
		fmt.Println(blockchain.CreateTx(*to, *amount))
	} else {
		TxCreateUsage()
	}
}
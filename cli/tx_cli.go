package cli

import (
	"fmt"
	"flag"
	"letsgo/blockchain"
	"letsgo/util"
	"encoding/hex"
	"letsgo/script"
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
	fmt.Printf("    %d inputs:\n", len(tx.Inputs))
	for idx, in := range tx.Inputs {
		fmt.Printf("    === [%d] ===\n", idx)
		fmt.Printf("    PrevHash: %x\n", in.PrevTransactionHash)
		fmt.Printf("    Vout: %d\n", util.DecodeInt(in.Vout))
		fmt.Printf("    ScriptSig: %s\n\n", script.Script.String(in.ScriptSig))
	}
	fmt.Printf("    %d outputs:\n", len(tx.Outputs))
	for idx, out := range tx.Outputs {
		fmt.Printf("    === [%d] ===\n", idx)
		fmt.Printf("    Value: %d\n", util.DecodeInt(out.Value))
		fmt.Printf("    ScriptPubKey: %s\n\n", script.Script.String(out.ScriptPubKey))
	}
}

func printTxOnly(tx *blockchain.Transaction){
	fmt.Printf("== TX %x ==\n", tx.GetHash())
	fmt.Printf("    Coinbase: %t\n", tx.IsCoinbase())
	fmt.Printf("    Version: %x\n", tx.Version)
	fmt.Printf("    Value %d\n\n", tx.GetValue())
	fmt.Printf("    %d inputs:\n", len(tx.Inputs))
	for idx, in := range tx.Inputs {
		fmt.Printf("    === [%d] ===\n", idx)
		fmt.Printf("    PrevHash: %x\n", in.PrevTransactionHash)
		fmt.Printf("    Vout: %d\n", util.DecodeInt(in.Vout))
		fmt.Printf("    ScriptSig: %s\n\n", script.Script.String(in.ScriptSig))
	}
	fmt.Printf("    %d outputs:\n", len(tx.Outputs))
	for idx, out := range tx.Outputs {
		fmt.Printf("    === [%d] ===\n", idx)
		fmt.Printf("    Value: %d\n", util.DecodeInt(out.Value))
		fmt.Printf("    ScriptPubKey: %s\n\n", script.Script.String(out.ScriptPubKey))
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
	from := TxCMD.String("from", "", "sender address")
	amount := TxCMD.Int("amount", 0, "amount to send")
	fees := TxCMD.Int("fees", 0, "fees to offer to miner")
	handleParsingError(TxCMD)

	if *to != "" && *amount > 0 {
		tx := blockchain.CreateTx(*from, *to, *amount, *fees)
		if tx != nil {
			block := blockchain.NewBlock([]blockchain.Transaction{*tx}, blockchain.BC.Tip)
			err := blockchain.BC.AddBlock(block)
			if err != nil {
				fmt.Println("Block non miné")
			}
		}
	} else {
		TxCreateUsage()
	}
}
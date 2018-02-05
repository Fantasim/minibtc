package cli

import (
	"fmt"
	"flag"
	"tway/blockchain"
	"tway/util"
	"tway/script"
)

func blockchainPrintUsage(){
	fmt.Println(" Options:")
	fmt.Println(" --basic \t Basic print with : <height> <hash> <size> <nb_tx>")
	fmt.Println(" --intermediate \t Intermediate print with <height> <hash> <size> <nb_tx> <merkle_root> <unix_time> <difficulty>")
	fmt.Println(" --advanced \t Advanced print with all block informations and tx details")
}


func printBasic(){
	e := blockchain.NewExplorer()
	i := blockchain.BC.Height
	for i > 0 {
		block := e.Next()
		fmt.Printf("============================== Block [%d] =============================\n", i - 1)
		fmt.Printf("Hash: %x\n", block.GetHash())
		fmt.Printf("Size: %d\n", util.DecodeInt(block.Size))
		fmt.Printf("Txs: %d\n\n", len(block.Transactions))
		i--
	}
}

func printIntermediate(){
	e := blockchain.NewExplorer()
	i := blockchain.BC.Height
	for i > 0 {
		block := e.Next()
		fmt.Printf("============================== Block [%d] =============================\n", i - 1)
		fmt.Printf("Hash: %x\n", block.GetHash())
		fmt.Printf("Merkle root: %x\n", block.Header.HashMerkleRoot)
		fmt.Printf("Size: %d\n", util.DecodeInt(block.Size))
		fmt.Printf("Unix time: %d\n", util.DecodeInt(block.Header.Time))
		fmt.Printf("Difficulty: %d\n", util.DecodeInt(block.Header.Bits))
		fmt.Printf("Txs: %d\n\n", len(block.Transactions))
		i--
	}
}

func printAdvanced(){
	e := blockchain.NewExplorer()
	i := blockchain.BC.Height
	for i > 0 {
		block := e.Next()
		
		fmt.Printf("============================== Block [%d] =============================\n", i - 1)
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
			if tx.IsCoinbase() == false {
				fmt.Println()
				for idx, in := range tx.Inputs {
					fmt.Printf("\t inputs[%d] Value: %s\n", idx, script.Script.String(in.ScriptSig))
				}
				fmt.Println()
			}
			for idx, out := range tx.Outputs {
				fmt.Printf("\t output[%d] Value: %d\n", idx, util.DecodeInt(out.Value))
				fmt.Printf("\t output[%d] scriptPubKey: %s\n", idx, script.Script.String(out.ScriptPubKey))
			}
		}
		i--
	}
}

func BlockchainPrintCli(){
	blockchainPrintCMD := flag.NewFlagSet("blockchain_print", flag.ExitOnError)
	basic := blockchainPrintCMD.Bool("basic", false, "Print blockchain with basic contents")
	intermediate := blockchainPrintCMD.Bool("intermediate", false, "Print blockchain with intermediate contents")
	advanced := blockchainPrintCMD.Bool("advanced", false, "Print blockchain with advanced contents")

	handleParsingError(blockchainPrintCMD)

	if *advanced == true {
		printAdvanced()
	} else if *intermediate == true {
		printIntermediate()
	} else if *basic == true {
		printBasic()
	}  else {
		blockchainPrintUsage()
	}
}

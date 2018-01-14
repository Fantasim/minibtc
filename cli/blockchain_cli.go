package cli

import (
	"fmt"
	"flag"
	"letsgo/blockchain"
	"strconv"
)

func blockchainUsage(){
	fmt.Println("blockchain")
	fmt.Println(" Options:")
	fmt.Println(" --remove \t Remove blockchain\n")
	fmt.Println("Others cmds starting by blockchain :")
	fmt.Println("\t blockchain print")
}

func blockchainPrintUsage(){
	fmt.Println("blockchain print")
	fmt.Println(" Options:")
	fmt.Println(" --basic \t Basic print with : <height> <hash> <size> <nb_tx>")
	fmt.Println(" --intermediate \t Intermediate print with <height> <hash> <size> <nb_tx> <merkle_root> <unix_time> <difficulty>")
	fmt.Println(" --advanced \t Advanced print with all block informations and tx details")
}

func PrintChain(option string){
	switch (option){
		case "basic":
			printBasic()
		default:
			return
	}
}

func printBasic(){
	e := blockchain.NewExplorer()
	i := blockchain.BC_HEIGHT
	for i > 0 {
		block := e.Next()
		fmt.Printf("============================== Block [%d] =============================\n", i - 1)
		fmt.Printf("Hash: %x\n", block.GetHash())
		size, _:= strconv.Atoi(string(block.Size))
		fmt.Printf("Size: %d\n", size)
		fmt.Printf("Txs: %d\n\n", len(block.Transactions))
		i--
	}
}

func printIntermediate(){
	e := blockchain.NewExplorer()
	i := blockchain.BC_HEIGHT
	for i > 0 {
		block := e.Next()
		fmt.Printf("============================== Block [%d] =============================\n", i - 1)
		fmt.Printf("Hash: %x\n", block.GetHash())
		fmt.Printf("Merkle root: %x\n", block.Header.HashMerkleRoot)
		size, _:= strconv.Atoi(string(block.Size))
		fmt.Printf("Size: %d\n", size)
		unix, _:= strconv.Atoi(string(block.Header.Time))
		fmt.Printf("Unix time: %d\n", unix)
		difficulty, _:= strconv.Atoi(string(block.Header.Bits))
		fmt.Printf("Difficulty: %d\n", difficulty )
		fmt.Printf("Txs: %d\n\n", len(block.Transactions))
		i--
	}
}

func printAdvanced(){
	e := blockchain.NewExplorer()
	i := blockchain.BC_HEIGHT
	for i > 0 {
		block := e.Next()
		
		fmt.Printf("============================== Block [%d] =============================\n", i - 1)
		fmt.Printf("Hash: %x\n", block.GetHash())
		fmt.Printf("Merkle root: %x\n", block.Header.HashMerkleRoot)
		size, _:= strconv.Atoi(string(block.Size))
		fmt.Printf("Size: %d\n", size)
		unix, _:= strconv.Atoi(string(block.Header.Time))
		fmt.Printf("Unix time: %d\n", unix)
		difficulty, _:= strconv.Atoi(string(block.Header.Bits))
		fmt.Printf("Difficulty: %d\n", difficulty )
		fmt.Printf("Txs: %d\n", len(block.Transactions))
		for idx, tx := range block.Transactions {
			fmt.Printf("\t=== Tx [%d] ===\n", idx)
			fmt.Printf("\t Coinbase: %t\n", tx.IsCoinbase())
			fmt.Printf("\t Hash: %x\n", tx.GetHash())
			fmt.Printf("\t Value %d\n", tx.GetValue())
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

func blockchainCli(){
	blockchainCMD := flag.NewFlagSet("blockchain", flag.ExitOnError)
	remove := blockchainCMD.Bool("remove", false, "Remove current blockchain if exist")
	handleParsingError(blockchainCMD)

	if *remove == true {
		blockchain.RemoveBlockchainDB()
	} else {
		blockchainUsage()
	}
}
package cli

import (
	"fmt"
	"flag"
	"tway/blockchain"
)

func blockchainUsage(){
	fmt.Println(" Options:")
	fmt.Println(" --remove \t Remove blockchain\n")
	fmt.Println("Others cmds starting by blockchain :")
	fmt.Println("\t blockchain_print")
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
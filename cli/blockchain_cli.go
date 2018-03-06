package cli

import (
	"fmt"
	"flag"
	"tway/blockchain"
	"tway/util"
)

func blockchainUsage(){
	fmt.Println(" Options:")
	fmt.Println(" --remove \t Remove blockchain\n")
	fmt.Println(" --average-block-time \t Print the average time a block take to be mined")
	fmt.Println("Others cmds starting by blockchain :")
	fmt.Println("\t blockchain_print")
}

func blockchainCli(){
	blockchainCMD := flag.NewFlagSet("blockchain", flag.ExitOnError)
	remove := blockchainCMD.Bool("remove", false, "Remove current blockchain if exist")
	averageBlockTime := blockchainCMD.Bool("average-block-time", false, "Get the average mining time of a block")
	handleParsingError(blockchainCMD)

	if *remove == true {
		blockchain.RemoveBlockchainDB()
	} else if *averageBlockTime {
		lastBlock := blockchain.BC.GetLastBlock()
		genesisBlock := blockchain.BC.GetGenesisBlock()
		
		chainHeight := blockchain.BC.Height
		lastBlockTime := util.DecodeInt(lastBlock.Header.Time)
		genesisBlockTime := util.DecodeInt(genesisBlock.Header.Time)

		fmt.Println((lastBlockTime - genesisBlockTime) / chainHeight, "seconds")
	} else {
		blockchainUsage()
	}
}
package cli

import (
	"fmt"
	"os"
)

type CLI struct {
	//empty struct
}

func (cli *CLI) printUsage() {
	fmt.Println("Commands:")
	fmt.Println(" wallet \t Manage local wallets")
	fmt.Println(" blockchain \t Manage blockchain")
}

//Verifie les arguments
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

//Demarre le cli
func Start(){
	cli := new(CLI)
	cli.validateArgs()
	cli.listMenu()
}

//la liste des commandes
func (cli *CLI) listMenu(){

	switch os.Args[1] {
		case "wallet":
			walletCli()
		
		case "blockchain":
			blockchainCli()

		case "blockchain_print":
			BlockchainPrintCli()

		case "tx":
			TxPrintCli()

		case "tx_create":
			TxCreateCli()
		default: 
			cli.printUsage()
	}
}
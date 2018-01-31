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
	fmt.Println(" block \t Manage block")
	fmt.Println(" blockchain \t Manage blockchain")
	fmt.Println(" server \t Manage server")
	fmt.Println(" tx \t Manage transactions")
	fmt.Println(" wallet \t Manage local wallets")
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
		case "block":
			BlockPrintCli()

		case "blockchain":
			blockchainCli()

		case "blockchain_print":
			BlockchainPrintCli()

		case "server":
			serverCli()

		case "tx":
			TxPrintCli()

		case "tx_create":
			TxCreateCli()

		case "wallet":
			walletCli()
		default: 
			cli.printUsage()
	}
}
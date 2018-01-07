package cli

import (
	"flag"
	"letsgo/wallet"
	"fmt"
)

func walletUsage(){
	fmt.Println("wallet")
	fmt.Println(" Options:")
	fmt.Println("	--new		Generate a new wallet")
	fmt.Println("	--list		Print list of local wallets")
}

func walletCli(){
	walletCMD := flag.NewFlagSet("wallet", flag.ExitOnError)
	new := walletCMD.Bool("new", false, "Create a new wallet")
	list := walletCMD.Bool("list", false, "Print list of wallets stored")

	handleParsingError(walletCMD)

	if *list {
		//affiche la liste des addresses locals
		wallet.PrintAddressStored()
	} else if *new {
		//genere un nouveau wallet
		fmt.Println(wallet.GenerateWallet())
	} else {
		walletUsage()
	}
}

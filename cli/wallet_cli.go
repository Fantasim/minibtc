package cli

import (
	"flag"
	"letsgo/wallet"
	"letsgo/blockchain"
	"fmt"
	"github.com/bradfitz/slice"
	"bytes"
)

func walletUsage(){
	fmt.Println(" Options:")
	fmt.Println("	--new		Generate a new wallet")
	fmt.Println("	--list		Print list of local wallets")
}

//Afficher les adresses du wallet
func PrintAddressStored(){
	wsList := blockchain.Walletinfo.Ws

    slice.Sort(wsList[:], func(i, j int) bool {
        return bytes.Compare(wsList[i].Address, wsList[j].Address) < 0
    })

	for _, ws := range wsList {
		fmt.Println(string(ws.Wallet.GetAddress()), "\t", ws.Amount)
	}
}

func walletCli(){
	walletCMD := flag.NewFlagSet("wallet", flag.ExitOnError)
	new := walletCMD.Bool("new", false, "Create a new wallet")
	list := walletCMD.Bool("list", false, "Print list of wallets stored")

	handleParsingError(walletCMD)

	if *list {
		//affiche la liste des addresses locals
		PrintAddressStored()
	} else if *new {
		//genere un nouveau wallet
		fmt.Println(wallet.GenerateWallet())
	} else {
		walletUsage()
	}
}

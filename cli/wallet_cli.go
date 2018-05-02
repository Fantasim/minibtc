package cli

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"tway/wallet"

	"github.com/bradfitz/slice"
)

func walletUsage() {
	fmt.Println(" Options:")
	fmt.Println("	--new					Generate a new wallet")
	fmt.Println("	--list					Print list of local wallets")
	fmt.Println("	--total 				Print total amount available in local wallets")
	fmt.Println("	--pubkeyhash-to-addr 	Print addr from a public key hashed")
}

//Afficher les adresses du wallet
func PrintAddressStored(pubkey bool, privkey bool) {
	wsList := wallet.Walletinfo.Ws

	slice.Sort(wsList[:], func(i, j int) bool {
		return bytes.Compare(wsList[i].Address, wsList[j].Address) < 0
	})

	for _, ws := range wsList {
		fmt.Print(string(ws.W.GetAddress()), "\t", ws.Amount)
		if ws.AmountLockedByMultiSig != 0 {
			fmt.Print("\t", ws.AmountLockedByMultiSig)
		}
		if pubkey || privkey {
			fmt.Println()
		}
		if pubkey {
			fmt.Println("Public key:", hex.EncodeToString(ws.W.PublicKey))
		}
		if privkey {
			fmt.Println("Private key:", hex.EncodeToString(ws.W.PrivateKey.D.Bytes()))
		}
		fmt.Print("\n")
	}
}

func PrintTotalAmountAvailable() {
	wsList := wallet.Walletinfo.Ws
	var total int
	for _, ws := range wsList {
		total += ws.Amount
	}
	fmt.Println(total, "coins are free to spend")
}

func walletCli() {
	walletCMD := flag.NewFlagSet("wallet", flag.ExitOnError)
	new := walletCMD.Bool("new", false, "Create a new wallet")
	list := walletCMD.Bool("list", false, "Print list of wallets stored")
	total := walletCMD.Bool("total", false, "Print total amount available in wallets stored")
	pubkey := walletCMD.Bool("pubkey", false, "Print public key of each wallet stored. Only works with --list")
	privkey := walletCMD.Bool("privkey", false, "Print private key of each wallet stored. Only works with --list")
	pubkeyHToAddr := walletCMD.String("pubkeyhash-to-addr", "", "convert a pubKeyHash to an address")

	handleParsingError(walletCMD)

	if *pubkeyHToAddr != "" {
		pubKeyHashBytes, _ := hex.DecodeString(*pubkeyHToAddr)
		addr := wallet.GetAddressFromPubKeyHash(pubKeyHashBytes)
		fmt.Println(string(addr))
		return
	}
	if *list {
		//affiche la liste des addresses locals
		PrintAddressStored(*pubkey, *privkey)
	} else if *new {
		//genere un nouveau wallet
		w := wallet.NewWallet()
		addr := string(w.GetAddress())[:]
		wallet.WalletList[addr] = w
		wallet.SaveToFile()

		fmt.Println("address:", hex.EncodeToString(w.GetAddress()))
		fmt.Println("public key:", hex.EncodeToString(w.PublicKey))
	} else if *total {
		PrintTotalAmountAvailable()
	} else {
		walletUsage()
	}
}

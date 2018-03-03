package cli

import (
	"flag"
	"fmt"
	b "tway/blockchain"
	"tway/util"
	"bytes"
	"encoding/hex"
	"tway/wallet"
	conf "tway/config"
	"github.com/bradfitz/slice"
)

func utxoUsage(){
	fmt.Println(" Options:")
	fmt.Println("	--all		Print all UTXOs")
	fmt.Println("	--printTX	Print tx linked with each UTXO")
}


func printAll(printTX bool){
	UTXOs := b.BC.FindUTXO()

	for txid, outputs := range UTXOs {
		if printTX == true {
			txidBytes, _ := hex.DecodeString(txid)
			tx, _, height := b.GetTxByHash(txidBytes)
			fmt.Println("Block height:", height)
			printTxBasic(tx)
			fmt.Println()
		} else {
			fmt.Println("txID:", txid)
		}
		for _, output := range outputs.Outputs {
			fmt.Println("value:", util.DecodeInt(output.Value))
		}
		fmt.Println()
	}
}

func printMine(printTX bool){
	Walletinfo := wallet.Walletinfo

	_, localUTXO := Walletinfo.GetLocalUnspentOutputs(conf.MAX_COIN, "")		
	slice.Sort(localUTXO[:], func(i, j int) bool {
		return bytes.Compare(localUTXO[i].TxID, localUTXO[j].TxID) < 0
	})

	txIDPrinted := make(map[string]bool)

	for idx, localOutput := range localUTXO {
		txid := hex.EncodeToString(localOutput.TxID)
		if !txIDPrinted[txid] {
			if idx != 0 {
				fmt.Println()
			}
			if printTX == true {
				tx, _, height := b.GetTxByHash(localOutput.TxID)
				fmt.Println("Block height:", height)
				printTxBasic(tx)
				fmt.Println()
			} else {
				fmt.Println("txID:", txid)				
			}
			txIDPrinted[txid] = true
		}
		fmt.Println("value:",localOutput.Amount)
	}
}

func printLinkedWithTx(txID string, printTX bool){
	UTXOs := b.BC.FindUTXO()

	if _, ok := UTXOs[txID]; !ok {
		fmt.Println("any utxo for this tx")
		return
	}
	outputs := UTXOs[txID]

	if printTX == true {
		txidBytes, _ := hex.DecodeString(txID)
		tx, _, height := b.GetTxByHash(txidBytes)
		fmt.Println("Block height:", height)
		printTxBasic(tx)
		fmt.Println()
	}
	for _, output := range outputs.Outputs {
		fmt.Println("value:", util.DecodeInt(output.Value))
	}

}

func UTXOCli(){
	utxoCMD := flag.NewFlagSet("utxo", flag.ExitOnError)
	all := utxoCMD.Bool("all", false, "Create a new wallet")
	mine := utxoCMD.Bool("mine", false, "Print list of wallets stored")
	txid := utxoCMD.String("txid", "", "Print UTXO linked with a tx")
	printTX := utxoCMD.Bool("printTX", false, "Print tx linked with an utxo")

	handleParsingError(utxoCMD)
	if *all == true {
		printAll(*printTX)
	} else if *mine == true {
		printMine(*printTX)
	} else if *txid != "" {
		printLinkedWithTx(*txid, *printTX)
	} else {
		utxoUsage()
	}

}
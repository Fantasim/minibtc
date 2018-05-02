package cli

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	b "tway/blockchain"
	"tway/script"
	"tway/twayutil"
	"tway/util"
	"tway/wallet"
)

func TxPrintUsage() {
	fmt.Println(" Options:")
	fmt.Println(" --hash \t Print tx equal to hash")
	fmt.Println(" --sign \t	sign a transaction")
	fmt.Println(" --address \t select a wallet linked with this address")
	fmt.Println("Others cmds starting by tx :")
	fmt.Println("\t tx_reate")
}

func printTxBlockchain(tx *twayutil.Transaction, block *twayutil.Block, height int) {
	fmt.Printf("Block height: %d\n", height)
	fmt.Printf("Block hash: %x\n\n", block.GetHash())
	fmt.Printf("== TX %x ==\n", tx.GetHash())
	fmt.Printf("    Coinbase: %t\n", tx.IsCoinbase())
	fmt.Printf("    Version: %x\n", tx.Version)
	fmt.Printf("    Value %d\n\n", tx.GetValue())
	fmt.Printf("    %d inputs:\n", len(tx.Inputs))
	for idx, in := range tx.Inputs {
		fmt.Printf("    === [%d] ===\n", idx)
		fmt.Printf("    PrevHash: %x\n", in.PrevTransactionHash)
		fmt.Printf("    Vout: %d\n", util.DecodeInt(in.Vout))
		fmt.Printf("    ScriptSig: %s\n\n", script.Script.String(in.ScriptSig))
	}
	fmt.Printf("    %d outputs:\n", len(tx.Outputs))
	for idx, out := range tx.Outputs {
		fmt.Printf("    === [%d] ===\n", idx)
		fmt.Printf("    Value: %d\n", util.DecodeInt(out.Value))
		fmt.Printf("    ScriptPubKey: %s\n\n", script.Script.String(out.ScriptPubKey))
	}
}

func printTx(tx *twayutil.Transaction) {
	fmt.Printf("== TX %x ==\n", tx.GetHash())
	fmt.Printf("    Coinbase: %t\n", tx.IsCoinbase())
	fmt.Printf("    Version: %x\n", tx.Version)
	fmt.Printf("    Value %d\n\n", tx.GetValue())
	fmt.Printf("    %d inputs:\n", len(tx.Inputs))
	for idx, in := range tx.Inputs {
		fmt.Printf("    === [%d] ===\n", idx)
		fmt.Printf("    PrevHash: %x\n", in.PrevTransactionHash)
		fmt.Printf("    Vout: %d\n", util.DecodeInt(in.Vout))
		fmt.Printf("    ScriptSig: %s\n\n", script.Script.String(in.ScriptSig))
	}
	fmt.Printf("    %d outputs:\n", len(tx.Outputs))
	for idx, out := range tx.Outputs {
		fmt.Printf("    === [%d] ===\n", idx)
		fmt.Printf("    Value: %d\n", util.DecodeInt(out.Value))
		fmt.Printf("    ScriptPubKey: %s\n\n", script.Script.String(out.ScriptPubKey))
	}
}

func printTxBasic(tx *twayutil.Transaction) {
	fmt.Printf("== TX %x ==\n", tx.GetHash())
	fmt.Printf("    Coinbase: %t\n", tx.IsCoinbase())
	fmt.Printf("    Value %d\n\n", tx.GetValue())
}

func TxPrintCli() {
	TxCMD := flag.NewFlagSet("tx", flag.ExitOnError)
	hash := TxCMD.String("hash", "", "Print tx if exist")
	sign := TxCMD.String("sign", "", "Sign a transaction by its txid")
	address := TxCMD.String("address", "", "Select a wallet linked by address")
	handleParsingError(TxCMD)

	if *hash != "" {
		h, _ := hex.DecodeString(*hash)
		tx, block, height := b.GetTxByHash(h)
		if height != -1 {
			printTxBlockchain(tx, block, height)
		}
	} else if *sign != "" && *address != "" {
		h, _ := hex.DecodeString(*sign)
		tx, _, _ := b.GetTxByHash(h)
		w := wallet.WalletList[*address]
		r, s, err := ecdsa.Sign(rand.Reader, &w.PrivateKey, tx.ToTxUtil().Serialize())
		if err != nil {
			fmt.Println(err)
			log.Panic(err)
		}

		signature := append(r.Bytes(), s.Bytes()...)
		fmt.Println("signature:", hex.EncodeToString(signature))

	} else {
		TxPrintUsage()
	}
}

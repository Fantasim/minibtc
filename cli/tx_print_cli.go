package cli

import ( 
	"fmt"
	"tway/util"
	"tway/twayutil"
	"tway/script"
	b "tway/blockchain"
	"flag"
	"encoding/hex"
)

func TxCreateUsage(){
	fmt.Println(" Options:")
	fmt.Println(" --to \t address to send")
	fmt.Println(" --amount \t amount to send")
}

func printTxBlockchain(tx *twayutil.Transaction, block *twayutil.Block, height int){
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

func printTx(tx *twayutil.Transaction){
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


func TxPrintCli(){
	TxCMD := flag.NewFlagSet("tx", flag.ExitOnError)
	hash := TxCMD.String("hash", "", "Print tx if exist")
	handleParsingError(TxCMD)

	if *hash != "" {
		h, _ := hex.DecodeString(*hash)
		tx, block, height := b.GetTxByHash(h)
		if height != -1 {
			printTxBlockchain(tx, block, height)
		}
	} else {
		TxPrintUsage()
	}
}
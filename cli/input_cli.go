package cli

import (
	"encoding/hex"
	"flag"
	"fmt"
	"strings"
	"tway/script"
	"tway/twayutil"
	"tway/util"
)

func inputUsage() {
	fmt.Println(" Options:")
	fmt.Println(" --new \t Create a new input")
	fmt.Println(" --vout \t index of the output to spend in the previous transaction \t /!| works with --new")
	fmt.Println(" --prevTxHash \t hash of the previous transaction \t /!| works with --new")
	fmt.Println(" --scriptSig \t signature script \t /!| works with --new")
}

func inputCli() {
	inputCMD := flag.NewFlagSet("input", flag.ExitOnError)
	new := inputCMD.Bool("new", false, "Create an input")
	vout := inputCMD.Int("vout", -1, "index of the output to spend in the previous transaction")
	prevTxHash := inputCMD.String("prevTxHash", "", "hash of the previous transaction")
	scriptSigString := inputCMD.String("scriptSig", "", "signature script")
	handleParsingError(inputCMD)

	if *new == true {

		if *vout > -1 && *prevTxHash != "" && *scriptSigString != "" {
			prevTxHashBytes, _ := hex.DecodeString(*prevTxHash)
			voutBytes := util.EncodeInt(*vout)

			var scriptSig [][]byte
			opcodeArr := strings.Split(*scriptSigString, " ")
			for _, op := range opcodeArr {
				b, found := script.GetOpcodeValueByName(op)
				if found == true {
					scriptSig = append(scriptSig, []byte{b})
				} else {
					bytes, _ := hex.DecodeString(op)
					scriptSig = append(scriptSig, bytes)
				}
			}
			in := twayutil.NewTxInput(prevTxHashBytes, voutBytes, scriptSig)
			fmt.Println(hex.EncodeToString(in.Serialize()))
		}

	} else {
		inputUsage()
	}
}

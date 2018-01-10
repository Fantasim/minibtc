package main

import (
	"letsgo/cli"
	txscript "letsgo/script"
	"letsgo/util"
	"fmt"
)

func TestScript(){
	script := []int{txscript.OP_DATA_1, txscript.OP_DATA_2, txscript.OP_ADD, txscript.OP_DATA_4, txscript.OP_EQUALVERIFY}
	engine := txscript.NewEngine()
	err := engine.Run(util.IntArrayToByteDoubleArray(script))
	if err == nil {
		fmt.Println("Script correct")
	} else {
		fmt.Println("Script incorrect")
	}
}

func main(){
	cli.Start()
}
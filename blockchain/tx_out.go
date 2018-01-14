package blockchain

import (
	"letsgo/util"
	"strconv"
)

type Output struct {
	Value []byte //[1-8]
	TxScriptLength []byte //[1-9]
	ScriptPubKey [][]byte
}

func NewTxOutput(scriptPubKey [][]byte, value int) Output {
	txo := Output{
		Value: util.IntToArrayByte(value),
		TxScriptLength: []byte(strconv.Itoa(util.LenDoubleSliceByte(scriptPubKey))),
		ScriptPubKey: scriptPubKey,
	}
	return txo
}
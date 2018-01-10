package blockchain

import (
	"letsgo/util"
)

type Output struct {
	Value []byte //[1-8]
	TxScriptLength []byte //[1-9]
	ScriptPubKey []byte
}

func NewTxOutput(){
	txo := &Output{Value: util.IntToArrayByte(REWARD)}
}
package blockchain

import (
	"letsgo/util"
	"strconv"
)

type Input struct {
	PrevTransactionHash []byte //[32]
	Vout []byte //[4]
	TxInScriptLen []byte //[1-9]
	ScriptSig [][]byte 
}

func NewTxInput(PrevTransactionHash []byte, Vout []byte, ScriptSig [][]byte) Input {
	in := Input{
		PrevTransactionHash: PrevTransactionHash,
		Vout: Vout,
		TxInScriptLen: []byte(strconv.Itoa(util.LenDoubleSliceByte(ScriptSig))),
		ScriptSig: ScriptSig,
	}
	return in
}

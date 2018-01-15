package blockchain

import (
	"letsgo/util"
	"strconv"
)

type TxInputs struct {
	Inputs []Input
}

type Input struct {
	PrevTransactionHash []byte //[32]
	Vout []byte //[4]
	TxInScriptLen []byte //[1-9]
	ScriptSig [][]byte 
}

//Retourne un nouvel input de tx
func NewTxInput(PrevTransactionHash []byte, Vout []byte, ScriptSig [][]byte) Input {
	in := Input{
		PrevTransactionHash: PrevTransactionHash,
		Vout: util.EncodeInt(-1),
		TxInScriptLen: []byte(strconv.Itoa(util.LenDoubleSliceByte(ScriptSig))),
		ScriptSig: ScriptSig,
	}
	return in
}

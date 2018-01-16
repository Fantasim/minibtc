package blockchain

import (
	"letsgo/util"
	"encoding/hex"
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
		TxInScriptLen: util.EncodeInt(util.LenDoubleSliceByte(ScriptSig)),
		ScriptSig: ScriptSig,
	}
	return in
}

func (in *Input) GetPrevTx() (string, *Transaction) {
	tx, _, _ := GetTxByHash(in.PrevTransactionHash)
	return hex.EncodeToString(tx.GetHash()), tx
}
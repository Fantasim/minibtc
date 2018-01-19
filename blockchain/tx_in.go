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
func NewTxInput(prevTransactionHash []byte, vout []byte, scriptSig [][]byte) Input {
	in := Input{
		PrevTransactionHash: prevTransactionHash,
		Vout: vout,
		TxInScriptLen: util.EncodeInt(util.LenDoubleSliceByte(scriptSig)),
		ScriptSig: scriptSig,
	}
	return in
}

//Récupère la transaction précédente d'un input
//retourne le txID et un pointer vers la transaction
func (in *Input) GetPrevTx() (string, *Transaction) {
	tx, _, _ := GetTxByHash(in.PrevTransactionHash)
	return hex.EncodeToString(tx.GetHash()), tx
}
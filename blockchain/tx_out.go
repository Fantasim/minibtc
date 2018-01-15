package blockchain

import (
	"letsgo/util"
	"bytes"
	"encoding/gob"
	"log"
	"letsgo/script"
	"fmt"
)

type TxOutputs struct {
	Outputs []Output
}

type Output struct {
	Value []byte //[1-8]
	TxScriptLength []byte //[1-9]
	ScriptPubKey [][]byte
}

//Retourne un nouvel output de tx
func NewTxOutput(scriptPubKey [][]byte, value int) Output {
	fmt.Println(util.EncodeInt(value))
	txo := Output{
		Value: util.EncodeInt(value),
		TxScriptLength: util.EncodeInt(util.LenDoubleSliceByte(scriptPubKey)),
		ScriptPubKey: scriptPubKey,
	}
	return txo
}

//Si l'output a été locké avec pubKeyHash
func (output *Output) IsLockedWithPubKeyHash(pubKeyHash []byte) bool {
	//on génère un scriptPubKey de type Pay-to-PubkeyHash
	//avec la clé public hashée passée en paramètre
	scriptPubKey := script.Script.LockingScript(pubKeyHash)

	//si la longueur des deux scripts est différente
	if len(output.ScriptPubKey) != len(scriptPubKey) {
		return false
	}
	//pour chaque element du script
	for i := 0; i < len(scriptPubKey); i++ {
		//si l'element du scriptPubKey de l'output est différent
		//de l'element du scriptPubKey généré avec la pubKeyHash
		if bytes.Compare(scriptPubKey[i], output.ScriptPubKey[i]) != 0 {
			return false
		}
	}
	return true
}

//TxOutputs -> []byte
func (outs *TxOutputs) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(outs)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

//[]byte -> TxOutputs
func DeserializeTxOutputs(d []byte) *TxOutputs {
	var outs TxOutputs

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&outs)
	if err != nil {
		log.Panic(err)
	}
	return &outs
}
package blockchain

import (
	"log"
	"letsgo/script"
	"bytes"
	"encoding/gob"
	"letsgo/util"
	"strconv"
)

type Transaction struct {
	Version []byte //[4]
	InCounter []byte //[1-9]
	Inputs []Input
	OutCounter []byte //[1-9]
	Outputs []Output
	LockTime []byte //[4]
}

// Serialize returns a serialized Transaction
func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

 func NewCoinbaseTx(toPubKey, signature []byte) Transaction {
	txIn := NewTxInput([]byte{}, []byte("-1"), script.Script.CoinbaseLockingScript(toPubKey))
	txOut := NewTxOutput(script.Script.CoinbaseUnlockingScript(signature), REWARD)

	tx := Transaction{
		Version: []byte{VERSION},
		InCounter: []byte("1"),
		OutCounter: []byte("1"),
		LockTime: []byte{0},
	}
	tx.Inputs = []Input{txIn}
	tx.Outputs = []Output{txOut}
	return tx
}



// Retourne l'ID de la transaction
func (tx *Transaction) GetHash() []byte {
	return util.Sha256(tx.Serialize())
}

func (tx *Transaction) GetValue() int {
	val := 0
	for _, out := range tx.Outputs {
		outVal, _ := strconv.Atoi(string(out.Value)) 
		val += outVal
	}
	return val
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].PrevTransactionHash) == 0 && bytes.Compare(tx.Inputs[0].Vout, []byte("-1")) == 0
}
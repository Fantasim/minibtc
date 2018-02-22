package util

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log"
)

type Input struct {
	PrevTransactionHash []byte //[32]
	Vout []byte //[4]
	TxInScriptLen []byte //[1-9]
	ScriptSig [][]byte 
}

type Output struct {
	Value []byte //[1-8]
	TxScriptLength []byte //[1-9]
	ScriptPubKey [][]byte
}

type Transaction struct {
	Version []byte //[4]
	InCounter []byte //[1-9]
	Inputs []Input
	OutCounter []byte //[1-9]
	Outputs []Output
	LockTime []byte //[4]
}

//Transaction -> []byte
func (tx Transaction) Serialize() []byte {
	b, err := json.Marshal(tx)
	if err != nil {
        log.Panic(err)
	}
	bu := new(bytes.Buffer)
	enc := gob.NewEncoder(bu)
	err = enc.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return bu.Bytes()
}

// DeserializeTransaction deserializes a transaction
func DeserializeTransaction(data []byte) Transaction {
	var tx Transaction
	var dataByte []byte

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&dataByte)
	if err != nil {
		log.Panic(err)
	}
	json.Unmarshal(dataByte, &tx)
	return tx
}
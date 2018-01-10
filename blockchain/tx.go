package blockchain

import (
	"crypto/rand"
	"log"
	"fmt"
	"letsgo/util"
)

type Transaction struct {
	Version []byte //[4]
	InCounter []byte //[1-9]
	OutCounter []byte //[1-9]
	LockTime []byte //[4]
}

type Input struct {
	PrevTransactionHash []byte //[32]
	Vout []byte //[4]
	TxInScriptLen []byte //[1-9]
	ScriptSig []byte 
}

func NewCoinbaseTx(to, data string) Transaction {
	//si il n'y a pas de data
	if data == "" {
		randData := make([]byte, 20)
		_, err := rand.Read(randData)
		if err != nil {
			log.Panic(err)
		}
		//genere une random data
		data = fmt.Sprintf("%x", randData)
	}

	txIn := Input{[]byte{}, []byte(util.IntToHex(-1)), []byte(util.IntToHex(0)), []byte(util.IntToHex(0))}
}
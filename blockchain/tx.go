package blockchain

import (
	"log"
	"letsgo/script"
	"bytes"
	"encoding/gob"
	"letsgo/util"
)

type Transaction struct {
	Version []byte //[4]
	InCounter []byte //[1-9]
	Inputs []Input
	OutCounter []byte //[1-9]
	Outputs []Output
	LockTime []byte //[4]
}

//Serialise une transaction
func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

//Créer une transaction coinbase
func NewCoinbaseTx(toPubKey, signature []byte) Transaction {
	var empty [][]byte
	txIn := NewTxInput([]byte{}, util.EncodeInt(-1), empty)
	txOut := NewTxOutput(script.Script.LockingScript(util.Sha256(toPubKey)), REWARD)

	tx := Transaction{
		Version: []byte{VERSION},
		InCounter: util.EncodeInt(1),
		OutCounter: util.EncodeInt(1),
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

//Retourne la valeur total des outputs de la TX
func (tx *Transaction) GetValue() int {
	val := 0
	for _, out := range tx.Outputs {
		val += util.DecodeInt(out.Value)
	}
	return val
}

//Retourne true si la tx est coinbase
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].PrevTransactionHash) == 0 && bytes.Compare(tx.Inputs[0].Vout, util.EncodeInt(-1)) == 0
}

//Récupère une transaction par son hash, avec le block dans lequel
//se trouve la transaction, ainsi que la hauteur du block
func GetTxByHash(hash []byte) (*Transaction, *Block, int) {
	be := NewExplorer()
	var i = BC_HEIGHT
	for i > 0 {
		block := be.Next()
		for _, tx := range block.Transactions {
			if bytes.Compare(hash, tx.GetHash()) == 0 {
				return &tx, block, i
			}
		}
		i--
	}
	return nil, nil, -1
}

func CreateTx(to string, amount int) *Transaction {
	return nil
}
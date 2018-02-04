package wire

import (
	"tway/util"
	"tway/script"
 	conf "tway/config"
	"bytes"
	"encoding/gob"
	"crypto/ecdsa"
	"crypto/rand"
	"log"
	"fmt"
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

func (in *Input) GetSize() uint64 {
	return 0
}

type Output struct {
	Value []byte //[1-8]
	TxScriptLength []byte //[1-9]
	ScriptPubKey [][]byte
}

type TxOutputs struct {
	Outputs []Output
}

//Retourne un nouvel output de tx
func NewTxOutput(scriptPubKey [][]byte, value int) Output {
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

func (out *Output) GetSize() uint64 {
	return 0
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
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

//Créer une transaction coinbase
func NewCoinbaseTx(toPubKey []byte, fees int) Transaction {
	var empty [][]byte
	txIn := NewTxInput([]byte{}, util.EncodeInt(-1), empty)
	txOut := NewTxOutput(script.Script.LockingScript(util.Ripemd160(util.Sha256(toPubKey))), conf.REWARD + fees)

	tx := Transaction{
		Version: []byte{conf.VERSION},
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

func (tx *Transaction) GetSize() uint64 {
	return 0
}

//Signe une transaction avec le clé privé
func (tx *Transaction) Sign(prevTxs map[string]*Transaction, inputsPrivKey []ecdsa.PrivateKey, inputsPubKey [][]byte){
	//si la transaction est coinbase
	if tx.IsCoinbase(){
		return
	}
	for idx, in := range tx.Inputs {
		//on signe les données
		prevTxid := hex.EncodeToString(in.PrevTransactionHash)
		r, s, err := ecdsa.Sign(rand.Reader, &inputsPrivKey[idx], prevTxs[prevTxid].Serialize())
		if err != nil {
			fmt.Println(err)
			log.Panic(err)
		}

		signature := append(r.Bytes(), s.Bytes()...)
		//on update l'input avec un nouvel input identique 
		//mais comprenant le bon scriptSig
		tx.Inputs[idx] = NewTxInput(tx.Inputs[idx].PrevTransactionHash, tx.Inputs[idx].Vout, script.Script.UnlockingScript(signature, inputsPubKey[idx]))
	}
}

//[]Transaction -> [][]byte
func TransactionToByteDoubleArray(txs []Transaction) [][]byte {
	ret := make([][]byte, len(txs))
	for idx, tx := range txs {
		ret[idx] = tx.Serialize()
	}
	return ret
}

func (tx *Transaction) GetFees(prevTxs map[string]*Transaction) int {
	if tx.IsCoinbase() == true {
		return 0
	}
	var total_input = 0
	var total_output = 0

	for _, out := range tx.Outputs {
		total_output += util.DecodeInt(out.Value)
	}

	for _, in := range tx.Inputs {
		prev := prevTxs[hex.EncodeToString(in.PrevTransactionHash)]
		for _, out := range prev.Outputs {
			total_input += util.DecodeInt(out.Value)
		}
	}
	return total_input - total_output
}
package blockchain

import (
	"log"
	"crypto/ecdsa"
	"crypto/rand"
	"letsgo/script"
	"bytes"
	"encoding/gob"
	"letsgo/util"
	"os"
	"letsgo/wallet"
	"fmt"
)

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
	var inputs []Input
	var outputs []Output

	if amount > Walletinfo.Amount {
		log.Println("You don't have enough coin to perform this transaction.")
		os.Exit(-1)
	}

	toPubKeyHash := wallet.GetPubKeyHashFromAddress([]byte(to))
	amountGot, localUnspents := Walletinfo.GetLocalUnspentOutputs(amount)
	
	for _, localUs := range localUnspents {
		var emptyScript [][]byte
		input := NewTxInput(localUs.TxID, util.EncodeInt(localUs.Output), emptyScript)
		inputs = append(inputs, input)
	}

	out := NewTxOutput(script.Script.LockingScript(toPubKeyHash), amount)
	outputs = append(outputs, out)
	
	if amountGot > amount {
		//on utilise la clé public du dernier output ajouté à la liste
		fromPubKeyHash := util.Sha256(localUnspents[len(localUnspents) - 1].Wallet.PublicKey)
		exc := NewTxOutput(script.Script.LockingScript(fromPubKeyHash), amountGot - amount)
		outputs = append(outputs, exc)
	}
	
	tx := &Transaction{
		Version: []byte{VERSION},
		InCounter: util.EncodeInt(len(inputs)),
		Inputs: inputs,
		OutCounter: util.EncodeInt(len(outputs)),
		Outputs: outputs,
	}
	tx.Sign(localUnspents, toPubKeyHash)
	return tx
	
}


func (tx *Transaction) Sign(localUnspents []LocalUnspentOutput, to []byte){
	if tx.IsCoinbase(){
		return
	}

	prevTXs := make(map[string]*Transaction)
	for _, in := range tx.Inputs {
		prevTxID, prevTx := in.GetPrevTx()
		prevTXs[prevTxID] = prevTx
	}
	var txCopy *Transaction
	txCopy = tx
	
	for idx, _ := range txCopy.Inputs {
		//transaction to string
		dataToSign := fmt.Sprintf("%x\n", txCopy)

		r, s, err := ecdsa.Sign(rand.Reader, &localUnspents[idx].Wallet.PrivateKey, []byte(dataToSign))
		if err != nil {
			fmt.Println(err)
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Inputs[idx] = NewTxInput(tx.Inputs[idx].PrevTransactionHash, tx.Inputs[idx].Vout, script.Script.UnlockingScript(signature, to))
	}
}

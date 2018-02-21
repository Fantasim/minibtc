package blockchain

import (
	"tway/twayutil"
	"bytes"
	"tway/util"
	"encoding/hex"
)

//Récupère une transaction par son hash, avec le block dans lequel
//se trouve la transaction, ainsi que la hauteur du block
func GetTxByHash(hash []byte) (*twayutil.Transaction, *twayutil.Block, int) {
	be := NewExplorer()
	var i = BC.Height
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

func GetPrevTxs(tx *twayutil.Transaction) map[string]*twayutil.Transaction {
	prevTXs := make(map[string]*twayutil.Transaction)

	for _, in := range tx.Inputs {
		prevTx, _, _ := GetTxByHash(in.PrevTransactionHash)
		prevTXs[hex.EncodeToString(in.PrevTransactionHash)] = prevTx
	}
	return prevTXs
}

func GetAmountsOutput(tx *twayutil.Transaction) int {
	var total_outputs = 0
	for _, out := range tx.Outputs {
		//on ajoute le montant au montant total redistribué vers une adresse.
		total_outputs += util.DecodeInt(out.Value)
	}
	return total_outputs
}

func GetAmountsInput(tx *twayutil.Transaction) int {
	var total_inputs = 0

	if tx.IsCoinbase() {
		return 0
	}
	//Pour chaque input de la tx
	for _, in := range tx.Inputs {
		//on recupere la transaction précédante de l'input
		prevTx, _, _ := GetTxByHash(in.PrevTransactionHash)
		//on récupère l'output ayant permis la création de l'input
		out := prevTx.Outputs[(util.DecodeInt(in.Vout))]
		//on ajoute le montant au montant total assemblés par les inputs
		total_inputs += util.DecodeInt(out.Value)
	}
	return total_inputs
}

//Retourne les informations concernant les montants de la transaction
//présent dans les inputs ou outputs
//Cette fonction retourne :
//montant total des inputs, montant total des outputs, frais de transactions
func GetAmounts(tx *twayutil.Transaction) (int, int, int) {
	var total_inputs = GetAmountsInput(tx)
	var total_outputs = GetAmountsOutput(tx)

	if tx.IsCoinbase() {
		return 0, total_outputs, 0
	}

	return total_inputs, total_outputs, total_inputs - total_outputs
}


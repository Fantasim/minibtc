package blockchain

import (
	"tway/wire"
	"bytes"
	"tway/util"
	"encoding/hex"
)

//Récupère une transaction par son hash, avec le block dans lequel
//se trouve la transaction, ainsi que la hauteur du block
func GetTxByHash(hash []byte) (*wire.Transaction, *wire.Block, int) {
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

func GetPrevTxs(tx *wire.Transaction) map[string]*wire.Transaction {
	prevTXs := make(map[string]*wire.Transaction)
	for _, in := range tx.Inputs {
		prevTx, _, _ := GetTxByHash(in.PrevTransactionHash)
		prevTXs[hex.EncodeToString(in.PrevTransactionHash)] = prevTx
	}
	return prevTXs
}

//Retourne les informations concernant les montants de la transaction
//présent dans les inputs ou outputs
//Cette fonction retourne :
//montant total des inputs, montant total des outputs, frais de transactions
func GetAmounts(hash []byte) (int, int, int) {

	var total_inputs = 0
	var total_outputs = 0

	tx, _, _ := GetTxByHash(hash)

	if tx.IsCoinbase() == true {
		return 0,0,0
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
	//pour chaque output de la tx
	for _, out := range tx.Outputs {
		//on ajoute le montant au montant total redistribué vers une adresse.
		total_outputs += util.DecodeInt(out.Value)
	}
	return total_inputs, total_outputs, total_inputs - total_outputs
}

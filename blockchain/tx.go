package blockchain

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	s "tway/script"
	"tway/twayutil"
	"tway/util"
)

//Cette fonction verifie chaque input de la transaction
//execute le scriptSig de l'input avec le scriptPubKey de l'output lié (Tx précédente)
func CheckIfTxIsCorrect(tx *twayutil.Transaction) error {
	if tx.IsCoinbase() == true {
		return nil
	}

	if err := CheckIfTxPutsAreCorrect(tx); err != nil {
		return err
	}

	//on recupere la liste des transactions ayant permis
	// la creation des inputs de la tx recu
	prevTXs := GetPrevTxs(tx)
	//pour chaque inputs
	for idx, in := range tx.Inputs {
		prevHash := hex.EncodeToString(in.PrevTransactionHash)

		if err := CheckIfInputIsAnUTXO(&in, prevTXs[prevHash]); err != nil {
			return err
		}

		vout := util.DecodeInt(in.Vout)
		//on recupère le script pubkey de la tx precente lié a cet input
		scriptPubKey := prevTXs[prevHash].Outputs[vout].ScriptPubKey
		scriptSig := in.ScriptSig

		//ScriptSig + ScriptPubKey
		scriptToRun := append(scriptSig, scriptPubKey...)
		//si le script n'est pas de type PubKeyHash
		/*if s.Script.IsPayToPubKeyHash(scriptToRun) == false {
			return errors.New(WRONG_SCRIPT)
		}*/
		//pour des raisons de fonctionnalités avec pkg on convertit le type twayutil.Transaction en type util.Transaction
		prevTXsUtil := make(map[string]*util.Transaction)
		for hash, tx := range prevTXs {
			prevTXsUtil[hash] = tx.ToTxUtil()
		}
		engine := s.NewEngine(prevTXsUtil, tx.ToTxUtil(), idx)
		//on execute le script
		err := engine.Run(scriptToRun)
		if err != nil {
			return err
		}
		//si la stack du script apres son execution n'est pas egale a true
		if engine.IsScriptSucceed() == false {
			return errors.New(WRONG_SCRIPT)
		}
	}
	return nil
}

//Verifie que le montant total amassé par les inputs est egal
//au montant total des outputs + frais de transaction
func CheckIfTxPutsAreCorrect(tx *twayutil.Transaction) error {
	if tx.IsCoinbase() == false {
		total_inputs, total_outputs, fees := GetAmounts(tx)
		if total_inputs != (total_outputs + fees) {
			return errors.New(WRONG_BLOCK_PUTS_VALUE)
		}
	}
	return nil
}

//Cette fonction verifie que l'output lié à l'input est un UTXO
func CheckIfInputIsAnUTXO(in *twayutil.Input, prevTX *twayutil.Transaction) error {
	vout := util.DecodeInt(in.Vout)
	unspentOutput := UTXO.GetUnSpentOutputByVoutAndTxHash(vout, prevTX.GetHash())
	if unspentOutput == nil {
		return errors.New(NOT_FOUND)
	}
	return nil
}

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

//Récupère la liste des transactions ayant permis la création de la totalité
//des inputs présents dans la transaction
func GetPrevTxs(tx *twayutil.Transaction) map[string]*twayutil.Transaction {
	prevTXs := make(map[string]*twayutil.Transaction)

	for _, in := range tx.Inputs {
		prevTx, _, h := GetTxByHash(in.PrevTransactionHash)
		if h > -1 {
			prevTXs[hex.EncodeToString(in.PrevTransactionHash)] = prevTx
		} else {
			fmt.Println("error in GetPrevTxs")
		}
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
		fmt.Println()
		prevTx, _, h := GetTxByHash(in.PrevTransactionHash)
		if h == -1 {
			fmt.Println("ERROR IN GetAmountsInput")
			return 0
		}
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

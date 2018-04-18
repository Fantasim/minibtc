package cli

import (
	"crypto/ecdsa"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	b "tway/blockchain"
	conf "tway/config"
	"tway/script"
	"tway/server"
	"tway/twayutil"
	"tway/util"
	"tway/wallet"
)

func TxPrintUsage() {
	fmt.Println(" Options:")
	fmt.Println(" --hash \t Print tx equal to hash")
	fmt.Println("Others cmds starting by tx :")
	fmt.Println("\t tx_create")
}

func createTx(from string, to string, amount int, fees int) *twayutil.Transaction {
	var inputs []twayutil.Input
	var inputsPubKey [][]byte
	var inputsPrivKey []ecdsa.PrivateKey
	var outputs []twayutil.Output
	var localUnspents []wallet.LocalUnspentOutput
	var amountGot int

	Walletinfo := wallet.Walletinfo

	if wallet.IsAddressValid(to) == false {
		fmt.Println("Address to send is not a valid address")
		return nil
	}
	//On récupère la clé public hashée à partir de l'address
	//à qui on envoie
	toPubKeyHash := wallet.GetPubKeyHashFromAddress([]byte(to))

	if from == "" {
		//on récupère une liste d'output qui totalise le montant a envoyer
		//on récupère aussi amountGot, qui est le total de la somme de value des outputs
		//Cette variable est indispensable, car si la valeur total obtenu est supérieur
		//au montant d'envoie, on doit transferer l'excédant sur le wallet du créateur de la tx
		amountGot, localUnspents = Walletinfo.GetLocalUnspentOutputs(amount+fees, to)
	} else {
		if wallet.IsAddressValid(from) == false {
			fmt.Println("sender address is not a valid address")
			return nil
		}
		amountGot, localUnspents = wallet.GetLocalUnspentOutputsByPubKeyHash(wallet.GetPubKeyHashFromAddress([]byte(from)), amount+fees)
	}

	//Si le montant d'envoie est inférieur au total des wallets locaux
	if (from == "" && (amount+fees) > Walletinfo.Amount) || (from != "" && (amount+fees) > amountGot) {
		log.Println("You don't have enough coin to perform this transaction.")
		return nil
	}

	//Pour chaque output
	for _, localUs := range localUnspents {
		var emptyScript [][]byte
		//on génère un input à partir de l'output
		input := twayutil.NewTxInput(localUs.TxID, util.EncodeInt(localUs.Idx), emptyScript)
		//et on l'ajoute à la liste
		inputs = append(inputs, input)
		//on ajoute dans un tableau de string la clé publique correspondant
		//au wallet proprietaire de l'output permettant la création de
		//cette input.
		inputsPubKey = append(inputsPubKey, localUs.W.PublicKey)
		//on ajoute dans un tableau de clé privée la clé privée correspondant
		//au wallet proprietaire de l'output permettant la signature de
		//cette input.
		//Ce tableau de clé privée permettra de signer chaque input.
		inputsPrivKey = append(inputsPrivKey, localUs.W.PrivateKey)
	}

	//on génére l'output vers l'address de notre destinaire
	out := twayutil.NewTxOutput(script.Script.LockingScript(toPubKeyHash), amount)
	outputs = append(outputs, out)

	//Si le montant récupére par les wallets locaux est supérieur
	//au montant que l'on décide d'envoyer
	if amountGot > (amount + fees) {
		//on utilise la clé public du dernier output ajouté à la liste
		fromPubKeyHash := wallet.HashPubKey(localUnspents[len(localUnspents)-1].W.PublicKey)
		//on génére un output vers le dernier output de la liste d'utxo récupéré
		//et on envoie l'excédant
		exc := twayutil.NewTxOutput(script.Script.LockingScript(fromPubKeyHash), amountGot-(amount+fees))
		outputs = append(outputs, exc)
	}
	tx := &twayutil.Transaction{
		Version:    []byte{conf.VERSION},
		InCounter:  util.EncodeInt(len(inputs)),
		Inputs:     inputs,
		OutCounter: util.EncodeInt(len(outputs)),
		Outputs:    outputs,
	}

	prevTXs := make(map[string]*util.Transaction)
	//on récupère la liste des transactions précédant
	//la liste des inputs de la tx
	for _, in := range tx.Inputs {
		prevTx, _, _ := b.GetTxByHash(in.PrevTransactionHash)
		txid := hex.EncodeToString(prevTx.GetHash())
		prevTXs[txid] = prevTx.ToTxUtil()
	}
	//on signe la transaction
	tx.Sign(prevTXs, inputsPrivKey, inputsPubKey)
	return tx
}

func TxCreateCli() {
	TxCMD := flag.NewFlagSet("tx_create", flag.ExitOnError)
	to := TxCMD.String("to", "", "address to send")
	from := TxCMD.String("from", "", "sender address")
	amount := TxCMD.Int("amount", 0, "amount to send")
	fees := TxCMD.Int("fees", 0, "fees to offer to miner")
	broadcast := TxCMD.Bool("broadcast", false, "broadcast transaction to the main node")
	handleParsingError(TxCMD)

	if *to != "" && *amount > 0 {
		tx := createTx(*from, *to, *amount, *fees)
		if tx == nil {
			return
		}
		if *broadcast == false {
			NewBlock([]twayutil.Transaction{*tx}, *fees)
		} else {
			s := server.NewServer(false, false, false)
			s.SendTx(server.GetMainNode(), tx)
		}
	} else {
		TxCreateUsage()
	}
}

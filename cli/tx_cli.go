package cli

import (
	"crypto/ecdsa"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"strings"
	b "tway/blockchain"
	conf "tway/config"
	"tway/script"
	"tway/server"
	"tway/twayutil"
	"tway/util"
	"tway/wallet"
)

func TxCreateUsage() {
	fmt.Println(" Options:")
	fmt.Println(" --to \t address to send \t //!\\ To create a PayToScriptHash TX, separate pubkeys with a ,")
	fmt.Println(" --broadcast \t send the transaction to network's nodes")
	fmt.Println(" --nSig \t number of signature required on the number of pubkeys to spent the TX.")
	fmt.Println(" --fees \t number of coins gived to the minor.")
	fmt.Println(" --from \t get utxos to create the transaction from the address linked with this field.")
	fmt.Println(" --amount \t amount to send")
}

type createTxInfo struct {
	from    string
	to      [][]byte
	amount  int
	fees    int
	nSig    int
	outputs []twayutil.Output
	inputs  []twayutil.Input
}

func createTx(ctxInfo createTxInfo) *twayutil.Transaction {
	var inputs []twayutil.Input
	var inputsPubKey [][]byte
	var inputsPrivKey []ecdsa.PrivateKey
	var outputs []twayutil.Output
	var localUnspents []wallet.LocalUnspentOutput
	var amountGot int

	from := ctxInfo.from
	to := ctxInfo.to
	amount := ctxInfo.amount
	fees := ctxInfo.fees
	nSig := ctxInfo.nSig
	inputs = ctxInfo.inputs
	outputs = ctxInfo.outputs

	Walletinfo := wallet.Walletinfo

	if len(ctxInfo.inputs) == 0 {

		if from == "" {
			//on récupère une liste d'output qui totalise le montant a envoyer
			//on récupère aussi amountGot, qui est le total de la somme de value des outputs
			//Cette variable est indispensable, car si la valeur total obtenu est supérieur
			//au montant d'envoie, on doit transferer l'excédant sur le wallet du créateur de la tx
			var notAcceptedAddr []byte
			if len(to) == 1 {
				notAcceptedAddr = wallet.GetAddressFromPubKeyHash(to[0])
			}

			amountGot, localUnspents = Walletinfo.GetLocalUnspentOutputs(amount+fees, string(notAcceptedAddr))
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

		//Pour chaque output non dépensé
		for _, localUs := range localUnspents {
			if localUs.AmountLockedByMultiSig > 0 {
				continue
			}
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

	} else {
		for _, in := range inputs {
			uo := b.UTXO.GetUnSpentOutputByVoutAndTxHash(util.DecodeInt(in.Vout), in.PrevTransactionHash)
			if uo == nil {
				log.Println("Wrong inputs")
				return nil
			}
			amountGot += util.DecodeInt(uo.Output.Value)
		}
		if amountGot > amount {
			fees = amount - amountGot
		} else if amount > amountGot {
			log.Println("You don't have enough coin to perform this transaction.")
			return nil
		}
	}

	if len(ctxInfo.outputs) == 0 {
		//on génére l'output vers l'address de notre destinaire
		out := twayutil.NewTxOutput(script.Script.LockingScript(to, nSig), amount)
		outputs = append(outputs, out)
	}

	//Si le montant récupére par les wallets locaux est supérieur
	//au montant que l'on décide d'envoyer
	if amountGot > (amount+fees) && len(ctxInfo.inputs) == 0 && len(ctxInfo.outputs) == 0 {
		//on utilise la clé public du dernier output ajouté à la liste
		fromPubKeyHash := [][]byte{wallet.HashPubKey(localUnspents[len(localUnspents)-1].W.PublicKey)}
		//on génére un output vers le dernier output de la liste d'utxo récupéré
		//et on envoie l'excédant
		exc := twayutil.NewTxOutput(script.Script.LockingScript(fromPubKeyHash, nSig), amountGot-(amount+fees))
		outputs = append(outputs, exc)
	}

	tx := &twayutil.Transaction{
		Version:    []byte{conf.VERSION},
		InCounter:  util.EncodeInt(len(inputs)),
		Inputs:     inputs,
		OutCounter: util.EncodeInt(len(outputs)),
		Outputs:    outputs,
	}

	if len(ctxInfo.inputs) > 0 {
		return tx
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

	//1. Une adresse (PayToPubKeyHash)
	//Exemple : "1NcwUvhJumC7Xjutq5zLcBkxxHBmEsrsLj"
	//2. Une liste de clé publique (PayToScriptHash | Multisig)
	//Exemple : "49998f7ef43d8aee5a41601cad951a4243c2c4ff1af5a174a1e705fcdbbdd7a79201f5d8f4d260d153072869b91d8167502b9940f1416dbd593916a85b726939, e742b84ad64924e94bbb3948b2dfd068a161656dc73ca592d46843435d67af6cc1509a8b5eb753837d30f1435a7d9ae66fdd5b23b1c509049511bb777087772a"
	toString := TxCMD.String("to", "", "address to send")
	//Une addresse en particulier possédant les UTXOs pour créer la transaction
	from := TxCMD.String("from", "", "sender address")
	//Montant de la transaction
	amount := TxCMD.Int("amount", 0, "amount to send")
	//Frais de transaction allant au mineur du block contenant la transaction
	fees := TxCMD.Int("fees", 0, "fees to offer to miner")
	//Nombre de signature nécéssaire pour dépenser une transaction multisig
	nSig := TxCMD.Int("nsig", 0, "Number of signature required to spend a pay to script hash tx")
	//Liste d'inputs au format hexadecimal
	//Exemple : "INPUT_HEX_1 INPUT_HEX_2 INPUT_HEX_3"
	inputsString := TxCMD.String("inputs", "", "Inputs manually created at hex format.")
	//Si spécifié, la transaction est envoyé au noeud principal qui la relaiera ensuite a tout le réseau
	broadcast := TxCMD.Bool("broadcast", false, "broadcast transaction to the main node")
	handleParsingError(TxCMD)

	var txInputs []twayutil.Input
	if *inputsString != "" {
		//chaque code hexadecimal representant un input est séparé par un espace blanc
		inputs := strings.Split(*inputsString, " ")
		for _, inHex := range inputs {
			inBytes, err := hex.DecodeString(inHex)
			if err == nil {
				in := twayutil.DeserializeInput(inBytes)
				if len(in.PrevTransactionHash) == 0 || len(in.ScriptSig) == 0 {
					log.Println("An input looks wrong.")
					return
				}
				txInputs = append(txInputs, *in)
			}
		}
	}

	var to [][]byte
	//si il y a plusieurs clé publique
	if strings.Contains(*toString, ",") {
		*toString = strings.Replace(*toString, " ", "", -1)
		for _, pkString := range strings.Split(*toString, ",") {
			pkBytes, _ := hex.DecodeString(pkString)
			to = append(to, pkBytes)
		}
		if *nSig == 0 {
			fmt.Println("\n/!\\ You must use --nsig parameter")
			return
		}
		//si il y a une addresse
	} else if *toString != "" {
		to = append(to, wallet.GetPubKeyHashFromAddress([]byte(*toString)))
	}
	if len(to) > 0 && *amount > 0 {
		ctxInfo := createTxInfo{*from, to, *amount, *fees, *nSig, []twayutil.Output{}, txInputs}
		tx := createTx(ctxInfo)

		if tx == nil {
			return
		}
		//Affichage de la transaction
		printTx(tx)
		//on mine un nouveau block localement
		if *broadcast == false {
			NewBlock([]twayutil.Transaction{*tx}, *fees)
		} else {
			//on l'envoie au main node qui la diffusera ensuite a tout le reseau
			s := server.NewServer(false, false, false)
			s.SendTx(server.GetMainNode(), tx)
		}
	} else {
		TxCreateUsage()
	}
}

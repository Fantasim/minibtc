package blockchain

import (
	"github.com/boltdb/bolt"
	"os"
	"log"
	"fmt"
	"encoding/hex"
	"letsgo/util"
	"letsgo/wallet"
	"time"
)

const (
	//Récompense de la transaction coinbase à chaque nouveau block miné
	REWARD = 50000000
	//total supply de la coin
	MAX_COIN = 21000000000000
	VERSION = byte(0x00)
	BLOCK_BUCKET = "blocks"
)

var (
	//Identifiant du noeud
	NODE_ID string

	DB_FILE = "/Users/fantasim/go/src/letsgo/assets/db/"

	BC *Blockchain 
	//Hauteur courante de la blockchain
	BC_HEIGHT int
	//Liste des outputs non depensés liés aux wallets locaux
	UnSpents []SpendableOutput
	//Total des montants disponible sur chacun des wallets locaux
	AmountAvailable int
)

type Blockchain struct {
	Tip []byte
	DB *bolt.DB
}

//Charge les outputs non dépensés lié aux wallets locaux
func loadSpendableOutputs(){
	for wallet.WalletLoaded == false {
		time.Sleep(100 * time.Millisecond)
	}
	AmountAvailable, UnSpents = UTXO.FindAllSpendableOutputs()
}

func init(){
	//Si la variable d'environnement NODE_ID est bien set
	NODE_ID = os.Getenv("NODE_ID")
	if NODE_ID == "" {
		fmt.Printf("Vous devez créer une variable d'environnement correspondant à l'ID de votre noeud.\nExemple : `export NODE_ID=10000`\n\n")
		os.Exit(1)
	}
	DB_FILE += NODE_ID
	//si la db existe déjà, on charge la blockchain
	if dbExists() == true {
		loadDB()
	} else {
		//block genèse
		genesis := NewGenesisBlock()
		//sinon on créer une blockchain à partir d'un block genèse
		if err := CreateBlockchainDB(genesis); err != nil {
			log.Panic(err)
		}
	}
	//on charge la hauteur de la blockchain
	BC_HEIGHT = BC.getHeight()
	//On réindex les utxo
	UTXO.Reindex()
	//On récupère les outputs non dépensé du wallet
	loadSpendableOutputs()
}

//récupère la height de la blockchain
func (b *Blockchain) getHeight() int {
	be := NewExplorer()
	var i = 0
	for {
		bl := be.Next();
		if bl == nil {
			break
		}
		i++
	}
	return i
}

//Récupère la totalité des utxos de la chain
func (b *Blockchain) FindUTXO() map[string]TxOutputs {
	UTXO := make(map[string]TxOutputs)
	spentTXOs := make(map[string][]int)
	e := NewExplorer()
	
	for {
		block := e.Next()
		//si le block a dépassé le genèse
		if block == nil {
			break;
		}
		//Pour chaque tx du block
		for _, tx := range block.Transactions {

			txID := hex.EncodeToString(tx.GetHash())
			
			Outputs:
				//parcours la liste des outputs de la tx
				for idx, out := range tx.Outputs {
					//si l'output a été ajouté à la liste des transaction dépensé
					if spentTXOs[txID] != nil {
						//pour chaque outputs correspondant à la transaction, ayant été dépensé
						for _, spentOutIdx := range spentTXOs[txID] {
							//si l'output a déjà été ajouté à la liste des outputs depensé
							if spentOutIdx == idx {
								continue Outputs
							}
						}
					}
					outs := UTXO[txID]
					outs.Outputs = append(outs.Outputs, out)
					UTXO[txID] = outs
				}
				//si la transaction n'est pas coinbase
				if tx.IsCoinbase() == false {
					//pour chaque input de la tx
					for _, in := range tx.Inputs {
						//On récupère la transaction précédente 
						prevHash := hex.EncodeToString(in.PrevTransactionHash)
						//On ajoute l'output lié à cet input dans la liste des outputs depensés
						spentTXOs[prevHash] = append(spentTXOs[prevHash], util.DecodeInt(in.Vout))
					}
				}
		}
	}
	return UTXO
} 
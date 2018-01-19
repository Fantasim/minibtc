package blockchain

import (
	"github.com/boltdb/bolt"
	"os"
	"log"
	"fmt"
	"encoding/hex"
	"letsgo/util"
	"letsgo/wallet"
	"bytes"
	"errors"
)

const (
	//Récompense de la transaction coinbase à chaque nouveau block miné
	REWARD = 50000000
	//total supply de la coin
	MAX_COIN = 21000000000000
	//Version du client
	VERSION = byte(0x00)
	//Nom du bucket des blocks
	BLOCK_BUCKET = "blocks"
)

var (
	// du noeud
	NODE_ID string

	//Path vers le fichier DB
	DB_FILE = "/Users/fantasim/go/src/letsgo/assets/db/"

	BC *Blockchain 

	//Hauteur courante de la blockchain
	BC_HEIGHT int
	//Liste des outputs non depensés liés aux wallets locaux
	Walletinfo *WalletInfo
)

type Blockchain struct {
	Tip []byte
	DB *bolt.DB
}

//Charge les outputs non dépensés lié aux wallets locaux
func loadSpendableOutputs(){
	if wallet.WalletLoaded == true {
		Walletinfo = GetWalletInfo()
	} else {
		fmt.Println("Les informations du wallets ne se chargent pas.")
		os.Exit(0)
	}
}

func init(){
	var blockchn = false
	//Si la variable d'environnement NODE_ID est bien set
	NODE_ID = os.Getenv("NODE_ID")
	if NODE_ID == "" {
		fmt.Printf("Vous devez créer une variable d'environnement correspondant à l'ID de votre noeud.\nExemple : `export NODE_ID=10000`\n\n")
		os.Exit(1)
	}
	//Ajoute la string de la variable d'environnement au path du fichier DB
	DB_FILE += NODE_ID
	//si la db existe déjà, on charge la blockchain
	if dbExists() == true {
		//charge le fichier db
		if loadDB() == nil {
			blockchn = true
		}
	} else {
		if len(wallet.WalletList) > 0 {
			var address string
			for addr, _ := range wallet.WalletList {
				address = addr
				break
			}
			//block genèse
			genesis := NewGenesisBlock(address)
			//sinon on créer une blockchain à partir d'un block genèse
			if err := CreateBlockchainDB(genesis); err != nil {
				log.Panic(err)
			}
			blockchn = true
		}
	}
	if blockchn == true {
		//on charge la hauteur de la blockchain
		BC_HEIGHT = BC.getHeight()
		//On réindex les utxo
		UTXO.Reindex()
		//On récupère les outputs non dépensé du wallet
		loadSpendableOutputs()
	}
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
	utxo := make(map[string]TxOutputs)
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
					//si un output dans la même transaction a déjà été ajouté dans la liste des spents txos
					if spentTXOs[txID] != nil {
						//pour chaque outputs correspondant à la transaction, ayant été dépensé
						for _, spentOutIdx := range spentTXOs[txID] {
							//si l'output a déjà été ajouté à la liste des outputs depensé
							if spentOutIdx == idx {
								continue Outputs
							}
						}
					}
					outs := utxo[txID]
					outs.Outputs = append(outs.Outputs, out)
					utxo[txID] = outs
				}

				if tx.IsCoinbase() == false {
					for _, in := range tx.Inputs {
						inTxID := hex.EncodeToString(in.PrevTransactionHash)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], util.DecodeInt(in.Vout))
					}
				}
				/*
				//si la transaction n'est pas coinbase
				if tx.IsCoinbase() == false {
					//pour chaque input de la tx
					for _, in := range tx.Inputs {
						//On récupère la transaction précédente
						prevHash := hex.EncodeToString(in.PrevTransactionHash)
						//On ajoute l'output lié à cet input dans la liste des outputs depensés
						spentTXOs[prevHash] = append(spentTXOs[prevHash], util.DecodeInt(in.Vout))
					}
				}*/
		}
	}
	return utxo
}

func (b *Blockchain) AddBlock(block *Block) error {
	db := BC.DB
	blockHash := block.GetHash()

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCK_BUCKET))
		//recupere dans la db un block correspondant au hash du nouveau block
		blockInDb := b.Get(blockHash)
		//si il existe deja
		if blockInDb != nil {
			fmt.Println("Le block", hex.EncodeToString(blockHash), "existe deja")
			return nil
		}
		//recupere le hash du block ayant la plus hauteur hauteur
		lastHash := b.Get([]byte("l"))
		lastBlockData := b.Get(lastHash)
		lastBlock := DeserializeBlock(lastBlockData)
		lastBlockHash := lastBlock.GetHash()

		if bytes.Compare(block.Header.HashPrevBlock, lastBlockHash) != 0 {
			return errors.New("New block is not the tip's next block")
		}

		//ajoute le block dans la db
		err := b.Put(blockHash, block.Serialize())
		if err != nil {
			return err
		}
		err = b.Put([]byte("l"), blockHash)
		if err != nil {
			log.Panic(err)
		}
		BC.Tip = blockHash
		return nil
	})
	if err == nil {
		BC_HEIGHT += 1
		UTXO.Reindex()
		loadSpendableOutputs()
	}
	return err
}
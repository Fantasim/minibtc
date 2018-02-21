package blockchain

import (
	"github.com/boltdb/bolt"
	"tway/twayutil"
	"tway/util"
	"errors"
	"encoding/hex"
	"fmt"
	"bytes"
	"log"
)

const (
	//Nom du bucket des blocks
	BLOCK_BUCKET = "blocks"
)

var (
	BC *Blockchain
	GENESIS_PUBKEY = []byte{189, 208, 30, 89, 219, 197, 16, 58, 25, 114, 192, 26, 220, 144, 175, 157, 49, 159, 118, 140, 125, 205, 53, 177, 7, 217, 176, 2, 32, 103, 6, 158, 41, 70, 93, 47, 232, 197, 86, 128, 148, 98, 99, 151, 120, 33, 166, 193, 45, 123, 29, 252, 213, 142, 130, 88, 248, 152, 109, 119, 89, 243, 129, 88}
)

type Blockchain struct {
	Tip []byte
	DB *bolt.DB
	Height int
}

func init(){
	if dbExists() == true {
		//charge le fichier db
		loadDB()
	} else {
		genesis := GenesisBlock(GENESIS_PUBKEY)
		CreateBlockchainDB(genesis)
	}
	UTXO.Reindex()
}

//récupère la height de la blockchain
func (b *Blockchain) getHeight() {
	be := NewExplorer()
	var i = 0
	for {
		bl := be.Next();
		if bl == nil {
			break
		}
		i++
	}
	b.Height = i
}

//Ajoute un block à la blockchain
func (b *Blockchain) AddBlock(block *twayutil.Block) error {
	db := b.DB

	if block == nil {
		return errors.New("nil block")
	}
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
		lastBlock := twayutil.DeserializeBlock(lastBlockData)
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
		BC.Height += 1
		UTXO.Reindex()
	}
	return err
}

//Récupère la totalité des utxos de la chain
func (b *Blockchain) FindUTXO() map[string]twayutil.TxOutputs {
	utxo := make(map[string]twayutil.TxOutputs)
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
	return utxo
}
package blockchain

import (
	"github.com/boltdb/bolt"
	"tway/twayutil"
	"tway/util"
	"encoding/hex"
	"log"
)

const (
	UTXO_BUCKET = "chainstate"
)

var (
	UTXO *UTXOSet
)

type UTXOSet struct {
}

//Structure représentant les informations liés à un UTXO
type UnspentOutput struct {
	TxID []byte
	Output int
	Amount int
}

//Récupère une liste d'outputs non dépensé locké avec le pubKeyHash
//d'un montant supérieur ou égal au montant passé en paramètre
func (utxo *UTXOSet) GetUnspentOutputsByPubKeyHash(pubKeyHash []byte, amount int) (int, []UnspentOutput) {
	var unspentOutputs []UnspentOutput
	accumulated := 0
	db := BC.DB

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UTXO_BUCKET))
		c := b.Cursor()

		//Pour chaque transaction comportant des outputs non dépensés
		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := twayutil.DeserializeTxOutputs(v)
			
			//pour chaque output non dépnesé de la tx
			for outIdx, out := range outs.Outputs {
				//si l'output est locké avec la pubKeyHash passé en paramètre
				//et que le montant accumulé est inférieur au montant passé en paramètre
				
				if out.IsLockedWithPubKeyHash(pubKeyHash) == true && accumulated < amount{
					value := util.DecodeInt(out.Value)
					accumulated += value
					//on ajoute l'output à la liste des utxo
					usOutput := UnspentOutput{k, outIdx, value}
					unspentOutputs = append(unspentOutputs, usOutput)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return accumulated, unspentOutputs
}

//Reindex la liste des utxo dans le bucket des UTXOS
func (utxo *UTXOSet) Reindex() error {
	bucketName := []byte(UTXO_BUCKET)
	db := BC.DB
	UTXO := BC.FindUTXO()
	
	err := db.Update(func (tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound {
			return err
		}
		_, err = tx.CreateBucket(bucketName)
		if err != nil {
			return err
		}
		b := tx.Bucket(bucketName)
		for txID, outs := range UTXO {
			key, _ := hex.DecodeString(txID) 
			err = b.Put(key, outs.Serialize())
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

//Compte le nombre de transaction contenant des outputs non dépensés
func (utxo *UTXOSet) CountTx() int {
	bucketName := []byte(UTXO_BUCKET)
	db := BC.DB
	var i = 0
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			i++
		}
		return nil
	})
	return i
}

package blockchain

import (
	"os"
	"github.com/boltdb/bolt"
	conf "tway/config"
	"tway/twayutil"
)

func dbExists() bool {
	if _, err := os.Stat(conf.DB_FILE); os.IsNotExist(err) {
		return false
	}
	return true
}

//Charge la db de la blockchain si elle existe
func loadDB() error {
	var tip []byte
	db, err := bolt.Open(conf.DB_FILE, 0600, nil)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCK_BUCKET))
		tip = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		return err
	}
	BC = &Blockchain{tip, db, 0}
	BC.getHeight()
	return nil
}

//Supprime la db de la blockchain
func RemoveBlockchainDB() error {
	return os.Remove(conf.DB_FILE)
}


//Créer une nouvelle blockchain avec le block genese contenant une tx coinbase
func CreateBlockchainDB(genesis *twayutil.Block) error {
	db, err := bolt.Open(conf.DB_FILE, 0600, nil)
	if err != nil {
		return err
	}

	var tip []byte

	err = db.Update(func(tx *bolt.Tx) error {
		//creer le bucket pour les blocks
		b, err := tx.CreateBucket([]byte(BLOCK_BUCKET))
		if err != nil {
			return err
		}
		//hash le block genese
		hash := genesis.GetHash()
		//ajoute dans ce bucket le block genese
		err = b.Put(hash, genesis.Serialize())
		if err != nil {
			return err
		}
		//ajoute le hash du dernier block
		err = b.Put([]byte("l"), hash)
		if err != nil {
			return err
		}
		tip = hash
		return nil
	})
	if err != nil {
		return err
	}
	//set le tip et la DB dans la structure Blockchain
	BC = &Blockchain{tip, db, 1}
	return nil
}

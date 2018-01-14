package blockchain

import (
	"os"
	"github.com/boltdb/bolt"
)

func dbExists() bool {
	if _, err := os.Stat(DB_FILE); os.IsNotExist(err) {
		return false
	}
	return true
}

func loadDB() error {
	var tip []byte
	db, err := bolt.Open(DB_FILE, 0600, nil)
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
	BC = &Blockchain{tip, db}
	return nil
}

func RemoveBlockchainDB() error {
	return os.Remove(DB_FILE)
}

func CreateBlockchainDB(genesis Block) error {
	db, err := bolt.Open(DB_FILE, 0600, nil)
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
	BC = &Blockchain{tip, db}
	return nil
}
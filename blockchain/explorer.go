package blockchain

import (
	"log"

	"github.com/boltdb/bolt"
	"bytes"
)

//Structure utiliser pour parcourir les blocks de la chain
type BlockchainExplorer struct {
	CurrentHash []byte
	DB          *bolt.DB
}

func NewExplorer() *BlockchainExplorer {
	return &BlockchainExplorer{CurrentHash: BC.Tip, DB: BC.DB}
}

//Retourne le block suivant 
//commence par le block correspondant au tip
func (be *BlockchainExplorer) Next() *Block{
	//Si le block est genese
	if bytes.Compare(be.CurrentHash,GENESIS_BLOCK_PREVHASH) == 0 {
		return nil
	}

	var block *Block

	err := be.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCK_BUCKET))
		encodedBlock := b.Get(be.CurrentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	be.CurrentHash = block.Header.HashPrevBlock
	return block
}

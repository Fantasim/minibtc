package blockchain

import (
	"github.com/boltdb/bolt"
	"os"
	"log"
	"fmt"
)

const (
	REWARD = 50000000
	VERSION = byte(0x00)
)

var (
	NODE_ID string
	DB_FILE = "/Users/fantasim/go/src/letsgo/assets/db/"
	BLOCK_BUCKET = "blocks"
	BC *Blockchain 
	BC_HEIGHT int
)

type Blockchain struct {
	Tip []byte
	DB *bolt.DB
}

func init(){
	NODE_ID = os.Getenv("NODE_ID")
	if NODE_ID == "" {
		fmt.Printf("Vous devez créer une variable d'environnement correspondant à l'ID de votre noeud.\nExemple : `export NODE_ID=10000`\n\n")
		os.Exit(1)
	}
	DB_FILE += NODE_ID
	if dbExists() == true {
		loadDB()
	} else {
		genesis := NewGenesisBlock()
		if err := CreateBlockchainDB(genesis); err != nil {
			log.Panic(err)
		}
	}
	BC_HEIGHT = BC.getHeight()
}

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

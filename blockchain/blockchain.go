package blockchain

import (
	"github.com/boltdb/bolt"
	"os"
	"fmt"
)

const (
	REWARD = 50000000
)

var (
	NODE_ID string
	DB_FILE = "/Users/fantasim/go/src/letsgo/assets/db/"
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
}

func CreateBlockchain(address string){
	if dbExists(DB_FILE) {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}


}



package config

import (
	"os"
	"fmt"
)

var (
	//HashPrevBlock du block genèse
	GENESIS_BLOCK_PREVHASH = []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
	//Identifiant du noeud
	NODE_ID string
	//Path vers le fichier DB
	DB_FILE = "/Users/fantasim/go/src/tway/assets/db/"
	WALLET_FILE = "/Users/fantasim/go/src/tway/assets/wallet/"
)

const (
	//Récompense de la transaction coinbase à chaque nouveau block miné
	REWARD = 50000000
	//total supply de la coin
	MAX_COIN = 21000000000000
	//Version du client
	VERSION = byte(0x00)

	NEW_DIFFICULTY_EACH_N_BLOCK = 5
	TARGET_TIME_BETWEEN_TWO_BLOCKS = 20 //20 seconds 
)

func InitPKG(){
	NODE_ID = os.Getenv("NODE_ID")
	//Si la variable d'environnement NODE_ID est bien set
	if NODE_ID == "" {
		fmt.Printf("Vous devez créer une variable d'environnement correspondant à l'ID de votre noeud.\nExemple : `export NODE_ID=10000`\n\n")
		os.Exit(1)
	}
	DB_FILE += NODE_ID
	WALLET_FILE += NODE_ID
}
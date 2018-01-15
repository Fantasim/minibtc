package wallet

import (
	"crypto/ecdsa"
	"letsgo/util"
	"os"
	"fmt"
)

const (
	Version = byte(0x00)
	AddressChecksumLen = 4 //checksumlen du Bitcoin
)

var (
	WalletList map[string]*Wallet
	WalletLoaded = false
	NODE_ID string
	WALLET_FILE = "/Users/fantasim/go/src/letsgo/assets/dat/"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func init(){
	NODE_ID = os.Getenv("NODE_ID")
	if NODE_ID == "" {
		fmt.Printf("Vous devez créer une variable d'environnement correspondant à l'ID de votre noeud.\nExemple : `export NODE_ID=10000`\n\n")
		os.Exit(1)
	}
	WALLET_FILE += NODE_ID
	WalletList = make(map[string]*Wallet)
	LoadFromFile()
	WalletLoaded = true
}

//Gènere un nouveau wallet
//Ajoute le wallet dans le fichier de stockage wallet du noeud 
//PWD = WalletFile + NODE_ID.dat
func GenerateWallet() string {
	w := NewWallet()
	addr := string(w.GetAddress())[:]
	//ajoute le wallet a la liste des wallets
	WalletList[addr] = w
	//met a jour le fichier .dat
	SaveToFile()
	return addr
}

//Génère un nouveau wallet
func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}

	return &wallet
}

//Formate la clé publique en address (processus utilisé par le BTC)
func (w Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{Version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := util.Base58Encode(fullPayload)

	return address
}

func GetPubKeyFromAddress(addr string) []byte {
	return WalletList[addr].PublicKey
}

//Recupere le checksum d'une clé publique (processus utilisé par le BTC)
func checksum(payload []byte) []byte {
	//double sha256
	doubleSha := util.Sha256(util.Sha256(payload))
	return doubleSha[:AddressChecksumLen]
}

//Hash la clé publique (processus utilisé par le BTC)
func HashPubKey(pubKey []byte) []byte {
	return util.Ripemd160(util.Sha256(pubKey))
}


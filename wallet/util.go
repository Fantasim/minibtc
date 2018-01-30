package wallet

import (
	"crypto/elliptic"
	"crypto/rand"
	"crypto/ecdsa"
	"log"
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"os"
)

//Génère une clé de pair (privée, publique)
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}

//Sauvegarde la liste des wallets dans le fichier .dat du client
func SaveToFile() {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(WalletList)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(WALLET_FILE, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

// LoadFromFile loads wallets from the file
func LoadFromFile() error {
	if _, err := os.Stat(WALLET_FILE); os.IsNotExist(err) {
		return err
	}

	fileContent, err := ioutil.ReadFile(WALLET_FILE)
	if err != nil {
		log.Panic(err)
	}

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&WalletList)
	if err != nil {
		log.Panic(err)
	}

	return nil
}
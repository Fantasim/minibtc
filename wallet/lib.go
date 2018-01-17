package wallet

import (
	"letsgo/util"
	"bytes"
	"errors"
	"crypto/rand"
	"crypto/ecdsa"
	mathr "math/rand"
)

//Vérifie qu'une adresse est correcte (processus utilisé par le BTC)
func IsAddressValid(addr string) bool {
	//base58 to pubkey hash
	pubKeyHash := util.Base58Decode([]byte(addr))
	//on recupere le checksum de la clé publique hashé
	actualChecksum := pubKeyHash[len(pubKeyHash)-AddressChecksumLen:]
	//on recupere la version
	version := pubKeyHash[0]
	//on recupere le contenu de la clé public hashé entre la version (Index = 1) et le checksum (Index = len - 4)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-AddressChecksumLen]

	//on créer un checksum correspondant au resultat de la reconstution de la clé public hashé
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	//si les deux checksum sont identiques l'adresse est valide
	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

func IsAddressStored(addr string) bool {
	return WalletList[addr] != nil
}

//Signe une data vide à partir de la clé privée correspondant
//a la clé publique sauvegarder dans le wallet
//retourne la signature
func SignPrivateKey(addr string) ([]byte, error) {
	if IsAddressStored(addr) == false {
		return []byte{}, errors.New("public key doesn't match with a private key stored")
	}
	
	w := *WalletList[addr]
	
	r, s, err := ecdsa.Sign(rand.Reader, &w.PrivateKey, []byte{})
	if err != nil {
		return []byte{}, err
	}
	signature := append(r.Bytes(), s.Bytes()...)
	return signature, nil
}

//Retourne un wallet aleatoire parmis les wallets locaux
func RandomWallet() *Wallet {
	random := mathr.Intn(len(WalletList) - 1)
	var i = 0
	for _, w := range WalletList {
		if i == random {
			return w
		}
		i++
	}
	return nil
}
package wallet

import (
	"letsgo/util"
	"bytes"
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
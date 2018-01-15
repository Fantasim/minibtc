package wallet

import (
	"fmt"
)

//Afficher les adresses du wallet
func PrintAddressStored(){
	for addr, _ := range WalletList {
		fmt.Println(addr)
	}
}


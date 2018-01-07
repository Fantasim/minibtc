package wallet

import (
	"fmt"
)

func PrintAddressStored(){
	for addr, _ := range WalletList {
		fmt.Println(addr)
	}
}


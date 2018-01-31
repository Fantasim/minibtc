package cli

import (
	"tway/server"
	"flag"
	"tway/wallet"
	"fmt"
)

func serverCli(){
	serverCMD := flag.NewFlagSet("server", flag.ExitOnError)
	minerAddress := serverCMD.String("miner", "", "Address to send reward of each block mined")
	handleParsingError(serverCMD)
	if *minerAddress != "" && wallet.IsAddressValid(*minerAddress) == false {
		fmt.Println("Miner address is not correct")
		return
	}
	server.StartServer(*minerAddress)
}
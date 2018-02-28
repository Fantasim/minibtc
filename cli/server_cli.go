package cli

import (
	"tway/server"
	"flag"
	"tway/wallet"
)

func serverCli(){
	serverCMD := flag.NewFlagSet("server", flag.ExitOnError)
	mining := serverCMD.Bool("mining", false, "enable mining")
	log := serverCMD.Bool("log", false, "Print logs")
	handleParsingError(serverCMD)
	if *mining == true && wallet.IsAWalletExist() == false {
		wallet.GenerateWallet()
	}
	s := server.NewServer(*log, *mining)
	s.StartServer()
}


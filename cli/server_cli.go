package cli

import (
	"flag"
	"fmt"

	"tway/server"
)

func ServerUsage() {
	fmt.Println(" Options:")
	fmt.Println(" --mining \t Enable mining")
	fmt.Println(" --log-server \t Print server's logs")
	fmt.Println(" --log-mining \t Print mining's logs")
}

func serverCli() {
	serverCMD := flag.NewFlagSet("server", flag.ExitOnError)

	mining := serverCMD.Bool("mining", false, "enable mining")
	logServer := serverCMD.Bool("log-server", false, "Print logs")
	logMining := serverCMD.Bool("log-mining", false, "Print mining logs")
	help := serverCMD.Bool("help", false, "Print usage of server CMD")

	handleParsingError(serverCMD)

	if *help == true {
		ServerUsage()
		return
	}

	s := server.NewServer(*logServer, *mining, *logMining)
	s.StartServer()
}

package cli

import (
	"flag"
	"os"
	"log"
)

func handleParsingError(set *flag.FlagSet) {
	err := set.Parse(os.Args[2:])
	if err != nil {
		log.Panic(err)
		os.Exit(2)
	}
}
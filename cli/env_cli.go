package cli

import (
	"fmt"
	"flag"
)

func EnvUsage(){
	fmt.Println(" Options:")
	fmt.Println(" --add \t add an env variable")
	fmt.Println(" --list \t list all env variable")
}

func EnvCli(){
	EnvCMD := flag.NewFlagSet("env", flag.ExitOnError)
	add := EnvCMD.String("add", "", "Add env variable")
	list := EnvCMD.Bool("list", false, "list environnement variables")
	handleParsingError(EnvCMD)

	if *add != "" {
		//
	} else if *list == true {
		//
	} else {
		EnvUsage()
	}
}
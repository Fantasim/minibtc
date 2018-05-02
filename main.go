package main

import (
	"math/rand"
	"time"
	"tway/blockchain"
	"tway/cli"
	"tway/config"
	"tway/wallet"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	config.InitPKG()
	blockchain.InitPKG()
	wallet.InitPKG()
}

func main() {
	cli.Start()
}

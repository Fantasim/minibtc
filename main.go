package main

import (
       "tway/cli"
       "math/rand"
       "time"
       "tway/wallet"
       "tway/blockchain"
       "tway/config"
)

func init(){
     rand.Seed(time.Now().UTC().UnixNano())
     config.InitPKG()
     blockchain.InitPKG()
     wallet.InitPKG()
}

func main(){
    cli.Start()
}
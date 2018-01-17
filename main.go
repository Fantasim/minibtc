package main

import (
	"letsgo/cli"
	"math/rand"
	"time"
)

func main(){
	rand.Seed(time.Now().UTC().UnixNano())
	cli.Start()
}

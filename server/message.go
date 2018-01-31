package server

import (
	"io/ioutil"
	"net"
	"log"
	"fmt"
	conf "tway/config"
)

var (
	KnownNodes []*NetAddress
	Node *NetAddress
)

//function appelé lorsqu'une nouvelle connexion est detectée
func HandleConnexion(conn net.Conn) {
	//on recupere le []byte dans request
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}

	command := bytesToCommand(request[:conf.CommandLength])
	switch command {
/*	case "addr":
		handleAddr(request)
	case "block":
		handleBlock(request)
	case "inv":
		handleInv(request, bc)
	case "getblocks":
		handleGetBlocks(request, bc)
	case "getdata":
		handleGetData(request, bc)
	case "tx":
		handleTx(request, bc)*/
	case "version":
		handleVersion(request)
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}
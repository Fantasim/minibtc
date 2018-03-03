package server

import (
	"io/ioutil"
	"net"
	"log"
	"fmt"
	conf "tway/config"
)

//function appelé lorsqu'une nouvelle connexion est detectée
func (s *Server) HandleConnexion(conn net.Conn) {
	//on recupere le []byte dans request
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := bytesToCommand(request[:conf.CommandLength])

	//5ms waiting in development env
	s.LocalWaiting()

	switch command {
	case "addr": //reception d'une liste d'adresse
		s.handleAddr(request)
	case "block": //reception d'un block
		s.handleBlock(request)
	case "inv": //reception d'une liste de hash (block ou transaction)
		s.handleInv(request)
	case "getaddr": //reception d'une demande de partage d'addresse
		s.handleAskAddr(request)
	case "getblocks": //reception d'une demande d'envoie de list de hash de block selon un intervalle de height donné
		s.handleAskBlocks(request)
	case "getdata": //reception d'une demande d'envoie de data grace au hash (tx ou block)
		s.handleGetData(request)
	case "ping": //reception d'un ping
		s.handlePing(request)
	case "pong": //reception d'une reponse à un ping
		s.handlePong(request)
/*	case "tx":
		handleTx(request, bc)*/
	case "verack": //reception d'une confirmation de reception de version
		s.handleVerack(request)
	case "version": //reception d'une version d'un noeud
		s.handleVersion(request)
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}
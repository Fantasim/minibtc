package server

import (
	"net"
	conf "tway/config"
	"log"
	"fmt"
)

func init(){
	Node = GetLocalNetAddr()
	KnownNodes = append(KnownNodes, GetMainNode())
}

func GetMainNode() *NetAddress {
	return NewNetAddressIPPort(conf.MainNodeIP, conf.MainNodePort)
}

//Demarrer le serveur du node
func StartServer(minerAddress string) {
	ln, err := net.Listen(conf.Protocol, Node.String())
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	//si l'adresse du noeud n'est pas un node connu
	if Node.IsEqual(KnownNodes[0]) == false  {
		//on envoie notre version de la blockchain au noeud principale
		sendVersion(KnownNodes[0])
	}
	fmt.Println("Main node:", Node.IsEqual(KnownNodes[0]) == true)
	fmt.Println("Running on", Node.String(), "\n")
	for {
		//attend le prochain appel
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go HandleConnexion(conn)
	}
	
}
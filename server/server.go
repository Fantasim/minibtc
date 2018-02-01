package server

import (
	"net"
	conf "tway/config"
	b "tway/blockchain"
	"log"
	"fmt"
	"sync"
	"io"
	"bytes"
)


type Server struct {
	version 				int32

	//ip of user who run server
	ipStatus				*NetAddress
	chain			 		*b.Blockchain

	mu					 	sync.Mutex
	peers             		map[string]*serverPeer 
}

func NewServer() *Server {
	s := &Server{
		version: conf.NodeVersion,
		ipStatus: GetLocalNetAddr(),
		peers: make(map[string]*serverPeer),
		chain: b.BC,
	}
	return s
}

func (s *Server) AddPeer(sp *serverPeer){
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.peers[sp.GetAddr()] == nil {
		s.peers[sp.GetAddr()] = sp
	}
}

func (s *Server) RemovePeer(sp *serverPeer){
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.peers, sp.GetAddr())
}

func (s *Server) sendData(addr string, data []byte) error {
	conn, err := net.Dial(conf.Protocol, addr)
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		
	} 
	defer conn.Close()

	//on envoie au noeud la data
	_, err = io.Copy(conn, bytes.NewReader(data))
	return err
}

//Demarrer le serveur du node
func (s *Server) StartServer(minerAddress string) {
	ln, err := net.Listen(conf.Protocol, s.ipStatus.String())
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	//si l'adresse du noeud n'est pas un node connu
	if s.ipStatus.IsEqual(GetMainNode()) == false  {
		//on envoie notre version de la blockchain au noeud principale
		s.sendVersion(GetMainNode())
	}
	fmt.Println("Main node:", s.ipStatus.IsEqual(GetMainNode()) == true)
	fmt.Println("Running on", s.ipStatus.String(), "\n")
	for {
		//attend le prochain appel
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go s.HandleConnexion(conn)
	}
	
}
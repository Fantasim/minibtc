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
	"time"
)



type Server struct {
	version 				int32

	log						bool
	//ip of user who run server
	ipStatus				*NetAddress
	chain			 		*b.Blockchain

	mu					 	sync.Mutex
	peers             		map[string]*serverPeer 
}

func NewServer(log bool) *Server {
	s := &Server{
		log: log,
		version: conf.NodeVersion,
		ipStatus: GetLocalNetAddr(),
		peers: make(map[string]*serverPeer),
		chain: b.BC,
	}
	return s
}

func (s *Server) Log(printTime bool, c... interface{}){
	if (s.log == true){
		if printTime == true {
			fmt.Print(time.Now().Format("15:04:05.000000"))
			fmt.Print(" ")
		}
		for _, seq := range c {
			fmt.Print(seq)
			fmt.Print(" ")
		}
		fmt.Print("\n")
	}
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

	fmt.Println("Main node:", s.ipStatus.IsEqual(GetMainNode()) == true)
	fmt.Println("Running on", s.ipStatus.String(), "\n")

	//si l'adresse du noeud n'est pas un node connu
	if s.ipStatus.IsEqual(GetMainNode()) == false  {

		addr := GetMainNode().String()
		s.AddPeer(NewServerPeer(addr))
		//on envoie notre version de la blockchain au noeud principale
		_, err := s.sendVersion(GetMainNode())
		if err == nil {
			s.peers[addr].VersionSent()
		}
	}
	
	for {
		//attend le prochain appel
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go s.HandleConnexion(conn)
	}
	
}
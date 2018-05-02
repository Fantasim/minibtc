package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
	b "tway/blockchain"
	conf "tway/config"
	mempool "tway/mempool"
	mine "tway/mining"
	peerhistory "tway/server/peerhistory"
	"tway/serverutil"
	"tway/twayutil"
)

type Server struct {
	version int32

	log    bool
	mining bool
	prod   bool
	//ip of user who run server
	ipStatus *serverutil.NetAddress
	chain    *b.Blockchain

	newFetchAtHeight int //when chain having this height, fetch next blocks to get best tip
	MiningManager    *mine.MiningManager
	BlockManager     *blockManager
	HistoryManager   *peerhistory.HistoryManager
	Mempool          *mempool.TxPool
	newBlock         chan *twayutil.Block
	mu               sync.Mutex
	addrMu           sync.Mutex
	peers            sync.Map
}

//Nouvelle structure Server
func NewServer(logServer bool, mining bool, logMining bool) *Server {
	s := &Server{
		log:            logServer,
		version:        conf.NodeVersion,
		ipStatus:       GetLocalNetAddr(),
		peers:          sync.Map{},
		MiningManager:  mine.NewMiningManager(b.BC.Tip, logMining, b.BC),
		BlockManager:   NewBlockManager(logServer, mining),
		HistoryManager: peerhistory.NewHistoryManager(true),
		Mempool:        mempool.Mempool,
		chain:          &*b.BC,
		mining:         mining,
		newBlock:       make(chan *twayutil.Block),
	}
	return s
}

func (s *Server) Log(printTime bool, c ...interface{}) {
	if s.log == true {
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

//Ajoute un nouveau pair
func (s *Server) AddPeer(sp *serverPeer) {
	sp.RequestReceived()
	s.peers.Store(sp.GetAddr(), sp)
}

func (s *Server) GetPeer(addr string) (*serverPeer, bool) {
	val, exist := s.peers.Load(addr)
	if exist == false {
		na, err := serverutil.NewNetAddressByString(addr)
		if err != nil {
			fmt.Println("ERROR IN GetPeer from serverutil.NewNetAddressByString(addr)")
			return nil, exist
		}
		newP := NewServerPeer(na)
		s.peers.Store(addr, newP)
		return newP, exist
	}
	ret := val.(*serverPeer)
	return ret, exist
}

//supprime un pair par son adresse
func (s *Server) RemovePeer(sp *serverPeer) {
	s.peers.Delete(sp.GetAddr())
}

//Envoie une data par requete TCP
func (s *Server) sendData(addr string, data []byte) error {
	conn, err := net.Dial(conf.Protocol, addr)
	if err != nil {
		s.Log(true, fmt.Sprintf("%s is not available\n", addr))
		return err
	}
	defer conn.Close()
	go func() {
		p, _ := s.GetPeer(addr)
		p.IncreaseBytesSent(uint64(len(data)))
		s.AddPeer(p)
	}()
	//on envoie au noeud la data
	_, err = io.Copy(conn, bytes.NewReader(data))
	return err
}

//Demarrer le serveur du node
func (s *Server) StartServer() {
	ln, err := net.Listen(conf.Protocol, s.ipStatus.String())
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()
	fmt.Println("Running on", s.ipStatus.String())
	fmt.Println("Current chain height:", b.BC.Height)
	fmt.Println("Main node:", s.ipStatus.IsEqual(GetMainNode()) == true, "\n")

	if s.mining == true {
		go s.HandleNewBlockMined()
	}

	//si l'adresse du noeud n'est pas un node connu
	if s.ipStatus.IsEqual(GetMainNode()) == false {
		go func() {
			addr := GetMainNode().String()
			na, err := serverutil.NewNetAddressByString(addr)
			if err != nil {
				fmt.Println("error from NewNetAddressByString in StartServer")
				return
			}
			s.AddPeer(NewServerPeer(na))
			//on envoie notre version de la blockchain au noeud principale
			s.sendVersion(GetMainNode())
		}()
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

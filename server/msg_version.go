package server

import (
	"time"
	conf "tway/config"
	"log"
	"fmt"
)

type MsgVersion struct {
	// Version of the protocol the node is using.
	ProtocolVersion int32
	// Time the message was generated.  This is encoded as an int64 on the wire.
	Timestamp time.Time
	// Address of the remote peer.
	AddrReceiver *NetAddress
	// Address of the local peer.
	AddrSender *NetAddress
	// Last block seen by the generator of the version message.
	LastBlock int
}

func (s *Server) NewVersion(addrTo *NetAddress) *MsgVersion {
	return &MsgVersion{
		ProtocolVersion: conf.NodeVersion,
		Timestamp: time.Now(),
		AddrReceiver: addrTo,
		AddrSender: s.ipStatus,
		LastBlock: s.chain.Height,
	}
}

//Envoie une structure de la version de notre blockchain au noeud principal
func (s *Server) sendVersion(addrTo *NetAddress) ([]byte, error) {
	payload := gobEncode(*s.NewVersion(addrTo))
	request := append(commandToBytes("version"), payload...)
	return request, s.sendData(addrTo.String(), request)
}

//Recupère la version d'un noeud
func (s *Server) handleVersion(request []byte) {
	var payload MsgVersion
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}

	go func(){
		addr := payload.AddrSender.String()
		s.AddPeer(NewServerPeer(addr))
		p := s.peers[addr]
		if p != nil {
			p.VersionSent()
			p.IncreaseBytesSent(uint64(len(request)))
			p.SetLastBlock(int64(payload.LastBlock))
			p.SetStartingHeight(int64(payload.LastBlock))
			if request, err := s.sendVerack(payload.AddrSender); err == nil {
				p.IncreaseBytesReceived(uint64(len(request)))
				p.VerAckReceived()
			}
			if p.IsConfirmed() == true {
				fmt.Println("Connexion successfully with", p.GetAddr())
				fmt.Println("Last block:", p.GetLastBlock())
			}
			s.peers[addr] = p
		}
	}()
	
	/*
	fmt.Println("Version received from :", payload.AddrSender.String())
	fmt.Println("Block height:", payload.LastBlock, "\n")
	*/

	//recupere la hauteur du noeud envoyant une version
	foreignerBestHeight := payload.LastBlock

	//si le height courant est inférieur au height du noeud recepteur
	if s.chain.Height < foreignerBestHeight {
		//on lui envoie une demande des blocks qu'il a
		//sendGetBlocks(payload.AddrFrom)
	} else if s.chain.Height > foreignerBestHeight  {
		//on lui envoie notre version
		s.sendVersion(payload.AddrSender)
	}
}
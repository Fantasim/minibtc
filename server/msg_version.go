package server

import (
	"time"
	conf "tway/config"
	"tway/blockchain"
	"log"
)

type MsgVersion struct {
	// Version of the protocol the node is using.
	ProtocolVersion int32
	// Time the message was generated.  This is encoded as an int64 on the twayutil.
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
	
	addr := addrTo.String()
	s.Log(true, "Version sent to:", addrTo.String())
	payload := gobEncode(*s.NewVersion(addrTo))
	request := append(commandToBytes("version"), payload...)
	err := s.sendData(addrTo.String(), request)
	if err == nil {
		s.peers[addr].VersionSent()
	}
	return request, err
}

//Recupère la version d'un noeud
func (s *Server) handleVersion(request []byte) {
	var payload MsgVersion
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}

	s.Log(false, "\n")
	s.Log(true, "Version received from :", payload.AddrSender.String())
	s.Log(false, "\t - Block height:", payload.LastBlock)
	s.Log(false, "\t - Version:", payload.ProtocolVersion, "\n")

	if blockchain.BC.Height < payload.LastBlock {
		s.sendAskBlocks(payload.AddrSender, [2]int{blockchain.BC.Height + 1, payload.LastBlock})
		//lui demander des blocks
	}

	//établie les informations concernant le pair
	//envoie un verack et sa version si non fait.
	go func(){
		addr := payload.AddrSender.String()
		s.AddPeer(NewServerPeer(addr))
		p := s.peers[addr]

		p.SetLastBlock(int64(payload.LastBlock))
		p.SetStartingHeight(int64(payload.LastBlock))
		p.HasSentVersion()
		if _, err := s.sendVerack(payload.AddrSender); err == nil {
			p := s.peers[addr]
			if p.IsVersionSent() == false {
				s.sendVersion(payload.AddrSender);
			}
		}
		s.peers[addr] = p
	}()
}
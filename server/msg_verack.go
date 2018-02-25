package server

import (
	"log"
	"time"
	conf "tway/config"
)

type MsgVerack struct {
	// Address of the local peer.
	AddrSender *NetAddress
	// Address of the local peer.
	AddrReceiver *NetAddress
}

func (s *Server) NewVerack(addrTo *NetAddress) *MsgVerack {
	return &MsgVerack{s.ipStatus, addrTo}
}

//Envoie une requete verack
func (s *Server) sendVerack(addrTo *NetAddress) ([]byte, error) {
	s.Log(true, "VerAck sent to:", addrTo.String())
	payload := gobEncode(*s.NewVerack(addrTo))
	request := append(commandToBytes("verack"), payload...)
	return request, s.sendData(addrTo.String(), request)
}

//Receptionne une requete verack
//confirmation de bonne reception d'une requete version 
func (s *Server) handleVerack(request []byte) {
	var payload MsgVerack
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}
	addr := payload.AddrSender.String()
	s.Log(true, "VerAck received from :", addr)

	p := s.peers[addr]
	p.VerAckReceived()
	s.peers[addr] = p
	//si les echanges de version ont été realisé et que la derniere demande d'address avec ce noeud date de plus de conf.TimeInMinuteBetween2AskAddrWithASameNode
	if p.IsVerAckReceived() && p.IsVersionSent() && time.Now().Add(time.Minute * conf.TimeInMinuteBetween2AskAddrWithASameNode * -1).After(time.Unix(0, p.GetLastAddrGetTime())) {
		//on demande une liste d'addresse
		s.sendAskAddr(payload.AddrSender)
	}
}
package server

import (
	"log"
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

func (s *Server) sendVerack(addrTo *NetAddress) ([]byte, error) {
	s.Log(true, "VerAck sent to:", addrTo.String())
	payload := gobEncode(*s.NewVerack(addrTo))
	request := append(commandToBytes("verack"), payload...)
	return request, s.sendData(addrTo.String(), request)
}

//Recupère la version d'un noeud
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
	if p.IsVerAckReceived() && p.IsVersionSent() && p.HasHeSentVersion() {
		s.sendAskAddr(payload.AddrSender)
	}
}
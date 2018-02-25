package server

import (
	"log"
)

type MsgPing struct {
	// Address of the local peer.
	AddrSender *NetAddress
	// Address of the local peer.
	AddrReceiver *NetAddress
}

func (s *Server) NewPing(addrTo *NetAddress) *MsgPing {
	return &MsgPing{s.ipStatus, addrTo}
}

//Envoie une requete ping
func (s *Server) sendPing(addrTo *NetAddress) ([]byte, error) {
	addr := addrTo.String()

	s.Log(true, "Ping sent to:", addr)
	payload := gobEncode(*s.NewPing(addrTo))
	request := append(commandToBytes("ping"), payload...)
	err := s.sendData(addrTo.String(), request)
	if err == nil {
		s.peers[addr].PingSent()
	}
	return request, err
}

//Receptionne une requete ping
func (s *Server) handlePing(request []byte) {
	var payload MsgPing
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}
	addr := payload.AddrSender.String()
	s.Log(true, "Ping received from :", addr)
	s.sendPong(payload.AddrSender)
}
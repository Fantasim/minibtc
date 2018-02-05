package server

import (
	"log"
)

type MsgAskAddr struct {
	// Address of the local peer.
	AddrSender *NetAddress
	// Address of the local peer.
	AddrReceiver *NetAddress
}

func (s *Server) NewAskAddr(addrTo *NetAddress) *MsgAskAddr {
	return &MsgAskAddr{s.ipStatus, addrTo}
}

func (s *Server) sendAskAddr(addrTo *NetAddress) ([]byte, error) {
	addr := addrTo.String()
	s.Log(true, "GetAddr sent to:", addr)
	//assigne en []byte la structure getblocks
	payload := gobEncode(*s.NewAskAddr(addrTo))
	//on append la commande et le payload
	request := append(commandToBytes("getaddr"), payload...)
	err := s.sendData(addr, request)
	if err == nil {
		s.peers[addr].AskAddr()
	}
	return request, err
}

//Recupère la version d'un noeud
func (s *Server) handleAskAddr(request []byte) {
	var payload MsgAskAddr
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}
	addr := payload.AddrSender.String()
	s.Log(true, "GetAddr received from :", addr)
	s.sendAddr(payload.AddrSender)
}
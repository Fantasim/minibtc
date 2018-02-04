package server

import (
	"log"
	"time"
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
	s.Log(true, "GetAddr sent to:", addrTo.String())
	//assigne en []byte la structure getblocks
	payload := gobEncode(*s.NewAskAddr(addrTo))
	//on append la commande et le payload
	request := append(commandToBytes("getaddr"), payload...)
	return request, s.sendData(addrTo.String(), request)
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
	p := s.peers[addr]
	if time.Now().Add(time.Second * -1800).After(time.Unix(p.GetLastAddrGetTime(), 0)) {
		s.sendAskAddr(payload.AddrSender)
	} 
}
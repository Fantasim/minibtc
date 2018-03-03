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

//Envoie une demande de liste de d'adresses
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

//Cette fonction permet de receptionner une demande de liste d'adresse
func (s *Server) handleAskAddr(request []byte) {
	var payload MsgAskAddr
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}
	addr := payload.AddrSender.String()
	s.peers[addr].IncreaseBytesReceived(uint64(len(request)))
	s.Log(true, "GetAddr received from :", addr)
	//envoie une liste d'adresse au noeud à l'origine de la requete
	s.sendAddr(payload.AddrSender)
}
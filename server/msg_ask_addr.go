package server

import (
	"log"
	"tway/serverutil"
)

func (s *Server) NewAskAddr(addrTo *serverutil.NetAddress) *serverutil.MsgAskAddr {
	return &serverutil.MsgAskAddr{s.ipStatus, addrTo}
}

//Envoie une demande de liste de d'adresses
func (s *Server) sendAskAddr(addrTo *serverutil.NetAddress) ([]byte, error) {
	addr := addrTo.String()
	s.Log(true, "GetAddr sent to:", addr)
	//assigne en []byte la structure getblocks
	payload := gobEncode(*s.NewAskAddr(addrTo))
	//on append la commande et le payload
	request := append(commandToBytes("getaddr"), payload...)
	err := s.sendData(addr, request)
	if err == nil {
		p, _ := s.GetPeer(addr)
		p.AskAddr()
		s.AddPeer(p)
	}
	return request, err
}

//Cette fonction permet de receptionner une demande de liste d'adresse
func (s *Server) handleAskAddr(request []byte) {
	var payload serverutil.MsgAskAddr
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}
	addr := payload.AddrSender.String()
	p, _ := s.GetPeer(addr)
	p.IncreaseBytesReceived(uint64(len(request)))
	s.AddPeer(p)
	s.Log(true, "GetAddr received from :", addr)
	//envoie une liste d'adresse au noeud à l'origine de la requete
	s.sendAddr(payload.AddrSender)
}

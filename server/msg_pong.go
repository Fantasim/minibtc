package server

import (
	"log"
	"tway/serverutil"
)

func (s *Server) NewPong(addrTo *serverutil.NetAddress) *serverutil.MsgPong {
	return &serverutil.MsgPong{s.ipStatus, addrTo}
}

//envoie une requete pong (reponse à un ping)
func (s *Server) sendPong(addrTo *serverutil.NetAddress) ([]byte, error) {
	addr := addrTo.String()

	s.Log(true, "Pong sent to:", addr)
	payload := gobEncode(*s.NewPong(addrTo))
	request := append(commandToBytes("pong"), payload...)
	return request, s.sendData(addrTo.String(), request)
}

//Receptionne une requete pong (reponse d'un ping)
func (s *Server) handlePong(request []byte) {
	var payload serverutil.MsgPong
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}
	addr := payload.AddrSender.String()
	s.Log(true, "Pong received from :", addr)
	p, _ := s.GetPeer(addr)
	p.PongReceived()
	p.IncreaseBytesReceived(uint64(len(request)))
	s.AddPeer(p)
}

package server

import (
	"log"
	"tway/serverutil"
)

func (s *Server) NewPing(addrTo *serverutil.NetAddress) *serverutil.MsgPing {
	return &serverutil.MsgPing{s.ipStatus, addrTo}
}

//Envoie une requete ping
func (s *Server) sendPing(addrTo *serverutil.NetAddress) ([]byte, error) {
	addr := addrTo.String()

	s.Log(true, "Ping sent to:", addr)
	payload := gobEncode(*s.NewPing(addrTo))
	request := append(commandToBytes("ping"), payload...)
	err := s.sendData(addrTo.String(), request)
	if err == nil {
		p, _ := s.GetPeer(addr)
		p.PingSent()
		s.AddPeer(p)
	}
	return request, err
}

//Receptionne une requete ping
func (s *Server) handlePing(request []byte) {
	var payload serverutil.MsgPing
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}
	addr := payload.AddrSender.String()
	p, _ := s.GetPeer(addr)
	p.IncreaseBytesReceived(uint64(len(request)))
	s.AddPeer(p)
	s.Log(true, "Ping received from :", addr)
	s.sendPong(payload.AddrSender)

}

package server

import (
	"log"
)

type MsgAddr struct {
	// Address of the local peer.
	AddrSender *NetAddress
	// Address of the local peer.
	AddrReceiver *NetAddress

	AddrList [][]byte
}

func (s *Server) GetAddrList() [][]byte{
	var ret [][]byte

	for _, peer := range s.peers {
		ret = append(ret, []byte(peer.GetAddr()))
	}
	return ret
}

func (s *Server) NewMsgAddr(addrTo *NetAddress) *MsgAddr {
	return &MsgAddr{s.ipStatus, addrTo, s.GetAddrList()}
}

func (s *Server) sendAddr(addrTo *NetAddress) ([]byte, error) {
	s.Log(true, "Addr sent to:", addrTo.String())
	//assigne en []byte la structure getblocks
	payload := gobEncode(*s.NewMsgAddr(addrTo))
	//on append la commande et le payload
	request := append(commandToBytes("addr"), payload...)
	return request, s.sendData(addrTo.String(), request)
}

func (s *Server) handleAddr(request []byte) {
	var payload MsgAddr
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}
	addr := payload.AddrSender.String()
	s.Log(true, "Addr received from :", addr)
	s.Log(false, "-", len(payload.AddrList), "adresses reçus")
	
	p := s.peers[addr]
	p.GotAddr()
	s.peers[addr] = p
	
	var nbNewPeers = 0
	for _, addrBytes := range payload.AddrList {
		addrString := string(addrBytes)
		if s.peers[addrString] == nil && addrString != s.ipStatus.String() {
			s.AddPeer(NewServerPeer(addrString))
			nbNewPeers++
		}
	}
	s.Log(false, "-", nbNewPeers , "nouveaux pairs")
	unTreatedPeers := s.ListOfUntreatedPeers()
	go s.treatPeersAfterPong(unTreatedPeers)
	
}
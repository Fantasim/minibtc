package server

import (
	"log"
	"tway/serverutil"
)

//Récupère la liste des adresses de confiance avec qui le noeud courant est ou a été contact
func (s *Server) GetAddrList() [][]byte {
	var ret [][]byte

	s.peers.Range(func(key, v interface{}) bool {
		ret = append(ret, []byte(key.(string)))
		return true
	})
	return ret
}

func (s *Server) NewMsgAddr(addrTo *serverutil.NetAddress) *serverutil.MsgAddr {
	return &serverutil.MsgAddr{s.ipStatus, addrTo, s.GetAddrList()}
}

//Envoie une liste d'adresse
func (s *Server) sendAddr(addrTo *serverutil.NetAddress) ([]byte, error) {
	s.Log(true, "Addr sent to:", addrTo.String())
	//assigne en []byte la structure getblocks
	payload := gobEncode(*s.NewMsgAddr(addrTo))
	//on append la commande et le payload
	request := append(commandToBytes("addr"), payload...)
	return request, s.sendData(addrTo.String(), request)
}

//cette fonction est appelée lorsque l'on recoit
// une liste d'addresse de nouveaux pairs.
func (s *Server) handleAddr(request []byte) {
	var payload serverutil.MsgAddr
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}

	addr := payload.AddrSender.String()
	s.Log(true, "Addr received from :", addr)
	s.Log(false, "-", len(payload.AddrList), "adresses reçus")

	p, _ := s.GetPeer(addr)
	p.GotAddr()
	p.IncreaseBytesReceived(uint64(len(request)))
	s.AddPeer(p)

	var nbNewPeers = 0
	for _, addrBytes := range payload.AddrList {
		addrString := string(addrBytes)
		//pour chaque nouveau pair si il n'existe pas déjà
		if addrString != s.ipStatus.String() {
			_, exist := s.GetPeer(addrString)
			if exist == false {
				nbNewPeers++
			}
		}
	}
	s.Log(false, "-", nbNewPeers, "nouveaux pairs")
	//on récupère une liste des pairs pour lesquelles on n'a recu
	//ni la version, ni un pong
	unTreatedPeers := s.ListOfUntreatedPeers()
	//Cette fonction envoie un ping a chacune de ces adresses et les traites selon la reponse obtenus.
	go s.treatPeersAfterPong(unTreatedPeers)

}

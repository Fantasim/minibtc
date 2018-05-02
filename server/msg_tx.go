package server

import (
	"encoding/hex"
	"log"
	"tway/serverutil"
	"tway/twayutil"
)

func (s *Server) NewMsgTx(addrTo *serverutil.NetAddress, tx *twayutil.Transaction) *serverutil.MsgTx {
	return &serverutil.MsgTx{s.ipStatus, addrTo, tx}
}

//Envoie un block
func (s *Server) SendTx(addrTo *serverutil.NetAddress, tx *twayutil.Transaction) ([]byte, error) {
	s.Log(true, "tx "+hex.EncodeToString(tx.GetHash())+" sent to:", addrTo.String())
	//assigne en []byte la structure getblocks
	payload := gobEncode(*s.NewMsgTx(addrTo, tx))
	//on append la commande et le payload
	request := append(commandToBytes("tx"), payload...)
	return request, s.sendData(addrTo.String(), request)
}

//Récéptionne un block
func (s *Server) handleTx(request []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var payload serverutil.MsgTx
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}

	/*
		addr := payload.AddrSender.String()
		p, _ := s.GetPeer(addr)
		p.IncreaseBytesReceived(uint64(len(request)))
		s.AddPeer(p)
	*/

	err := s.Mempool.AddTx(payload.Tx)
	if err != nil {
		log.Println("There is in error in handleTx, tx cannot be added.")
		return
	}

	if s.ipStatus.IsEqual(GetMainNode()) {
		s.peers.Range(func(key, val interface{}) bool {
			p := val.(*serverPeer)
			if payload.AddrSender.IsEqual(p.GetNetAddress()) == false {
				s.SendTx(p.GetNetAddress(), payload.Tx)
			}
			return true
		})
	}

}

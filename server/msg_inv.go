package server

import (
	"log"
	"encoding/hex"
)

type MsgInv struct {
	// Address of the local peer.
	AddrSender *NetAddress
	// Address of the local peer.
	AddrReceiver *NetAddress
	Kind string // "transaction" || "block"
	List [][]byte
}


func (s *Server) rangeTxList(data [][]byte){
}

func (s *Server) rangeBlockList(addrTo *NetAddress, data [][]byte, toSP *serverPeer){
	for _, item := range data {
		if b, _ := s.chain.GetBlockByHash(item); b == nil {
			_, err := s.sendGetData(addrTo, item, "block")
			if err == nil {
				s.blockmanager.StartDownloadBlock(hex.EncodeToString(item), toSP)
			}
		}
	}
}

func (s *Server) NewMsgInv(addrTo *NetAddress, kind string, list [][]byte) *MsgInv{
	return &MsgInv{s.ipStatus, addrTo, kind, list}
}

func (s *Server) sendInv(addrTo *NetAddress, kind string, list [][]byte) ([]byte, error) {
	s.Log(true, "Inv kind:"+kind+ " sent to:", addrTo.String())
	//assigne en []byte la structure getblocks
	payload := gobEncode(*s.NewMsgInv(addrTo, kind, list))
	//on append la commande et le payload
	request := append(commandToBytes("inv"), payload...)
	return request, s.sendData(addrTo.String(), request)
}

func (s *Server) handleInv(request []byte){
	var payload MsgInv
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}
	s.Log(true , "Inv kind:"+payload.Kind+" received from :", payload.AddrSender.String())
	s.Log(false, "list of", len(payload.List), payload.Kind)
	if payload.Kind == "block" {
		s.rangeBlockList(payload.AddrSender, payload.List, s.peers[payload.AddrSender.String()])
	} else {
		s.rangeTxList(payload.List)
	}

}
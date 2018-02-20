package server

import (
	"log"
	"tway/twayutil"
)

type MsgAskBlocks struct {
	// Address of the local peer.
	Addr *NetAddress
	Range [2]int
}

func (s *Server) NewMsgAskBlock(rng [2]int) *MsgAskBlocks{
	return &MsgAskBlocks{s.ipStatus, rng}
}

func (s *Server) sendAskBlocks(addrTo *NetAddress, rng [2]int) ([]byte, error) {
	s.Log(true, "GetBlocks sent to:", addrTo.String())
	//assigne en []byte la structure getblocks
	payload := gobEncode(*s.NewMsgAskBlock(rng))
	//on append la commande et le payload
	request := append(commandToBytes("getblocks"), payload...)
	return request, s.sendData(addrTo.String(), request)
}

func (s *Server) handleAskBlocks(request []byte){
	var payload MsgAskBlocks
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}
	listBlock := s.chain.GetNBlocksNextToHeight(payload.Range[0] - 1, payload.Range[1] - payload.Range[0] + 1)
	s.Log(true , "GetBlocks received from :", payload.Addr.String())
	s.Log(false, "height:", payload.Range[0] - 1)
	s.Log(false, len(listBlock), "blocks found")

	listHash := twayutil.GetListBlocksHash(listBlock)
	s.sendInv(payload.Addr, "block", listHash)
}
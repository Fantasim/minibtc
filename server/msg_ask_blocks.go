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
	askBlock := s.NewMsgAskBlock(rng)
	//assigne en []byte la structure getblocks
	payload := gobEncode(*askBlock)
	//on append la commande et le payload
	request := append(commandToBytes("getblocks"), payload...)
	err := s.sendData(addrTo.String(), request)
	if err == nil {
		s.HistoryManager.NewGetBlocksHistory(askBlock, true)
	}
	return request, err
}

//Receptionne une demande de liste de hash de block dans un intervalle de height donné
//voir structure MsgAskBlocks
func (s *Server) handleAskBlocks(request []byte){
	var payload MsgAskBlocks
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}

	s.HistoryManager.NewGetBlocksHistory(&payload, false)

	//récupère une liste de block dans un intervalle donné
	listBlock := s.chain.GetNBlocksNextToHeight(payload.Range[0] - 1, payload.Range[1] - payload.Range[0] + 1)
	s.Log(true , "GetBlocks received from :", payload.Addr.String())
	s.Log(false, "height:", payload.Range[0] - 1)
	s.Log(false, len(listBlock), "blocks found")
	
	//recupère une liste de hash depuis une list de block
	listHash := twayutil.GetListBlocksHashFromMap(listBlock)
	s.sendInv(payload.Addr, "block", listHash)
}
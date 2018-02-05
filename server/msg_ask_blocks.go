package server

import (
	"log"
	"tway/wire"
	"fmt"
	"encoding/hex"
)

type MsgAskBlocks struct {
	// Address of the local peer.
	Addr *NetAddress
	Height int
}

func (s *Server) NewMsgAskBlock() *MsgAskBlocks{
	return &MsgAskBlocks{s.ipStatus, s.chain.Height}
}

func (s *Server) sendAskBlocks(addrTo *NetAddress) ([]byte, error) {
	//assigne en []byte la structure getblocks
	payload := gobEncode(*s.NewMsgAskBlock())
	//on append la commande et le payload
	request := append(commandToBytes("getblocks"), payload...)
	return request, s.sendData(addrTo.String(), request)
}

func (s *Server) handleAskBlocks(request []byte){
	var payload MsgAskBlocks
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}

	if s.log == true {
		fmt.Println("Ask blocks received from :", payload.Addr.String())
		fmt.Println("height:", payload.Height)
	}
	listBlock := s.chain.GetNBlocksNextToHeight(payload.Height)
	listHash := wire.GetListBlocksHash(listBlock)
	for height, h := range listHash {
		fmt.Println(height, hex.EncodeToString(h))
	}
}
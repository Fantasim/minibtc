package server

import (
	"tway/twayutil"
	"encoding/hex"
	"log"
)

type MsgBlock struct {
	// Address of the local peer.
	AddrSender *NetAddress
	// Address of the local peer.
	AddrReceiver *NetAddress

	Data []byte
}

func (s *Server) NewMsgBlock(addrTo *NetAddress, data []byte) *MsgBlock {
	return &MsgBlock{s.ipStatus, addrTo, data}
}

//Envoie un block 
func (s *Server) sendBlock(addrTo *NetAddress, block *twayutil.Block) ([]byte, error) {
	s.Log(true, "block "+ hex.EncodeToString(block.GetHash()) +" sent to:", addrTo.String())
	//assigne en []byte la structure getblocks
	payload := gobEncode(*s.NewMsgBlock(addrTo, block.Serialize()))
	//on append la commande et le payload
	request := append(commandToBytes("block"), payload...)
	return request, s.sendData(addrTo.String(), request)
}

//Récéptionne un block
func (s *Server) handleBlock(request []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var payload MsgBlock
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}
	addr := payload.AddrSender.String()
	block := twayutil.DeserializeBlock(payload.Data)
	if block != nil {
		s.Log(true, "block "+ hex.EncodeToString(block.GetHash()) +" received from :", addr)
	} else {
		s.Log(true, "wrong block received from :", addr)		
	}

	s.Log(false, "handleBlock: current tip:", hex.EncodeToString(s.chain.Tip))
	s.Log(false, "handleBlock: block prev hash:", hex.EncodeToString(block.Header.HashPrevBlock))
	s.Log(false, "handleBlock: block hash:", hex.EncodeToString(block.GetHash()))
	s.BlockManager.BlockDownloaded(block, s)
}
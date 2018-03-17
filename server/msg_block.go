package server

import (
	"encoding/hex"
	"log"
	"tway/serverutil"
	"tway/twayutil"
)

func (s *Server) NewMsgBlock(addrTo *serverutil.NetAddress, data []byte) *serverutil.MsgBlock {
	return &serverutil.MsgBlock{s.ipStatus, addrTo, data}
}

//Envoie un block
func (s *Server) sendBlock(addrTo *serverutil.NetAddress, block *twayutil.Block) ([]byte, error) {
	s.Log(true, "block "+hex.EncodeToString(block.GetHash())+" sent to:", addrTo.String())
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

	var payload serverutil.MsgBlock
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}
	addr := payload.AddrSender.String()
	p, _ := s.GetPeer(addr)
	p.IncreaseBytesReceived(uint64(len(request)))
	s.AddPeer(p)

	block := twayutil.DeserializeBlock(payload.Data)
	if block != nil {
		s.Log(true, "block "+hex.EncodeToString(block.GetHash())+" received from :", addr)
	} else {
		s.Log(true, "wrong block received from :", addr)
	}

	s.BlockManager.BlockDownloaded(block, s)
}

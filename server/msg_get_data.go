package server

import (
	"encoding/hex"
	"fmt"
	"log"
	"time"
	"tway/mempool"
	"tway/serverutil"
)

func (s *Server) NewMsgGetData(addrTo *serverutil.NetAddress, ID []byte, kind string) *serverutil.MsgGetData {
	return &serverutil.MsgGetData{s.ipStatus, addrTo, ID, kind}
}

func (s *Server) sendGetData(addrTo *serverutil.NetAddress, ID []byte, kind string) ([]byte, error) {
	s.Log(true, fmt.Sprintf("GetData kind: %s, with ID:%s sent to %s", kind, hex.EncodeToString(ID), addrTo.String()))
	//assigne en []byte la structure getblocks
	payload := gobEncode(*s.NewMsgGetData(addrTo, ID, kind))
	//on append la commande et le payload
	request := append(commandToBytes("getdata"), payload...)
	return request, s.sendData(addrTo.String(), request)
}

//Receptionne une demande de block ou de transaction
func (s *Server) handleGetData(request []byte) {
	var payload serverutil.MsgGetData
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}
	addr := payload.AddrSender.String()
	p, _ := s.GetPeer(addr)
	p.IncreaseBytesReceived(uint64(len(request)))
	s.AddPeer(p)
	s.Log(true, fmt.Sprintf("GetData kind: %s, with ID:%s received from %s", payload.Kind, hex.EncodeToString(payload.ID), addr))

	if payload.Kind == "block" {
		//block
		//on recupère le block si il existe
		block, _ := s.chain.GetBlockByHash(payload.ID)
		if block != nil {
			//envoie le block au noeud créateur de la requete
			s.sendBlock(payload.AddrSender, block)
		} else {
			fmt.Println("block is nil :( handleGetData")
			go func() {
				for {
					time.Sleep(time.Millisecond * 50)
					block, _ := s.chain.GetBlockByHash(payload.ID)
					if block != nil {
						s.sendBlock(payload.AddrSender, block)
					}
				}
			}()
		}
	} else {
		tx := mempool.Mempool.GetTx(hex.EncodeToString(payload.ID))
		if tx != nil {
			s.SendTx(payload.AddrSender, tx)
		}
	}
}

package server

import (
	"log"
	"fmt"
	"encoding/hex"
	"time"
)

type MsgGetData struct {
	// Address of the local peer.
	AddrSender *NetAddress
	// Address of the local peer.
	AddrReceiver *NetAddress

	ID 			[]byte //hash du block ou de la tx
	Kind 		string //"block" ou "tx"
}

func (s *Server) NewMsgGetData(addrTo *NetAddress, ID []byte, kind string) *MsgGetData {
	return &MsgGetData{s.ipStatus, addrTo, ID, kind}
}

func (s *Server) sendGetData(addrTo *NetAddress, ID []byte, kind string) ([]byte, error) {
	s.Log(true, fmt.Sprintf("GetData kind: %s, with ID:%s sent to %s", kind, hex.EncodeToString(ID), addrTo.String()))
	//assigne en []byte la structure getblocks
	payload := gobEncode(*s.NewMsgGetData(addrTo, ID, kind))
	//on append la commande et le payload
	request := append(commandToBytes("getdata"), payload...)
	return request, s.sendData(addrTo.String(), request)
}

//Receptionne une demande de block ou de transaction 
func (s *Server) handleGetData(request []byte) {
	var payload MsgGetData
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}
	addr := payload.AddrSender.String()
	s.peers[addr].IncreaseBytesReceived(uint64(len(request)))
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
			go func(){
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
		//tx
	}
}
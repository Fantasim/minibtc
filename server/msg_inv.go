package server

import (
	"log"
	"encoding/hex"
	"tway/util"
	conf "tway/config"
)

type MsgInv struct {
	// Address of the local peer.
	AddrSender *NetAddress
	// Address of the local peer.
	AddrReceiver *NetAddress
	Kind string // "tx" || "block"
	List [][]byte
}


func (s *Server) rangeTxList(data [][]byte){
}

//parcours une liste hash de block suite a une requete handleInv
func (s *Server) rangeBlockList(addrTo *NetAddress, data [][]byte, toSP *serverPeer, heightExpectedOfFirstElem int){
	for idx, item := range data {
		//on recupère le block correspondant au hash, si il existe
		if b, _ := s.chain.GetBlockByHash(item); b == nil {
			hashBlock := hex.EncodeToString(item)
			//si le block n'est pas ajouté dans la liste des blocks a téléchargé du block manager
			if s.BlockManager.download[hashBlock] == nil {
				//on demande le block au noeud qui a envoyé la liste de blocks
				_, err := s.sendGetData(addrTo, item, "block")
				if err == nil {
					//on indique au block manager que l'on commence a télécharger le block
					var heightExpected = -1
					if heightExpectedOfFirstElem != -1 {
						heightExpected = heightExpectedOfFirstElem + idx
					}
					s.BlockManager.StartDownloadBlock(hashBlock, toSP, int64(heightExpected))
				}
			}
		}
	}
}

func (s *Server) NewMsgInv(addrTo *NetAddress, kind string, list [][]byte) *MsgInv{
	return &MsgInv{s.ipStatus, addrTo, kind, list}
}

//Envoie une liste de hash de data (block || tx) 
func (s *Server) sendInv(addrTo *NetAddress, kind string, list [][]byte) ([]byte, error) {
	s.Log(true, "Inv kind:"+kind+ " sent to:", addrTo.String())
	//assigne en []byte la structure getblocks
	payload := gobEncode(*s.NewMsgInv(addrTo, kind, list))
	//on append la commande et le payload
	request := append(commandToBytes("inv"), payload...)
	return request, s.sendData(addrTo.String(), request)
}

//Retourne le pourcentage de succès des envoie de requêtes
func (s *Server) BootstrapInv(kind string, list [][]byte) float64 {
	var nbRequestSucceeded = 0
	peers := s.GetCloserAndSafestPeers()
	for addr, _ := range peers {
		na := NewNetAddressIPPort(util.StringToNetIpAndPort(addr))
		_, err := s.sendInv(na, kind, list)
		if err == nil {
			nbRequestSucceeded++
		}
	}
	return float64(len(peers)) / float64(nbRequestSucceeded) * 100
}


//Receptionne une liste de hash de data (block || tx)
func (s *Server) handleInv(request []byte){
	var payload MsgInv	
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}

	addr := payload.AddrSender.String()
	s.peers[addr].IncreaseBytesReceived(uint64(len(request)))
	s.Log(true , "Inv kind:"+payload.Kind+" received from :", addr)
	s.Log(false, "list of", len(payload.List), payload.Kind)

	if payload.Kind == "block" {
		var gbh *getBlocksHistory
		var heightExpectedOfFirstElem = -1

		gbh = s.HistoryManager.GetBlock[addr].Select(true).sortByDate(true).first()
		if gbh != nil {
			if gbh.Message.Range[1] - gbh.Message.Range[0] + 1 >= conf.MaxBlockPerMsg {
				heightExpectedOfFirstElem = gbh.Message.Range[0]
			}
		}
		s.rangeBlockList(payload.AddrSender, payload.List, s.peers[addr], heightExpectedOfFirstElem)		
	} else {
		s.rangeTxList(payload.List)
	}
}
package server

import (
	"encoding/hex"
	"log"
	conf "tway/config"
	"tway/mempool"
	"tway/server/peerhistory"
	"tway/serverutil"
	"tway/util"
)

func (s *Server) rangeTxList(addrTo *serverutil.NetAddress, data [][]byte) {

	var indexToAsk []int
	for idx, item := range data {
		if mempool.Mempool.StartDownloadTx(item) == nil {
			indexToAsk = append(indexToAsk, idx)
		}
	}

	for _, idx := range indexToAsk {
		_, err := s.sendGetData(addrTo, data[idx], "tx")
		if err != nil {
			mempool.Mempool.RemoveDownloadInformation(data[idx])
		}
	}
}

//parcours une liste hash de block suite a une requete handleInv
func (s *Server) rangeBlockList(addrTo *serverutil.NetAddress, data [][]byte, toSP *serverPeer, heightExpectedOfFirstElem int) {
	for idx, item := range data {
		//on recupère le block correspondant au hash, si il existe
		if b, _ := s.chain.GetBlockByHash(item); b == nil {
			hashBlock := hex.EncodeToString(item)
			//si le block n'est pas ajouté dans la liste des blocks a téléchargé du block manager
			if s.BlockManager.download[hashBlock] == nil {
				//on demande le block au noeud qui a envoyé la liste de blocks
				_, err := s.sendGetData(addrTo, item, "block")
				if err == nil {
					var heightExpected = -1
					if heightExpectedOfFirstElem != -1 {
						heightExpected = heightExpectedOfFirstElem + idx
					}
					//on indique au block manager que l'on commence a télécharger le block
					s.BlockManager.StartDownloadBlock(hashBlock, toSP, int64(heightExpected))
				}
			}
		}
	}
}

func (s *Server) NewMsgInv(addrTo *serverutil.NetAddress, kind string, list [][]byte) *serverutil.MsgInv {
	return &serverutil.MsgInv{s.ipStatus, addrTo, kind, list}
}

//Envoie une liste de hash de data (block || tx)
func (s *Server) sendInv(addrTo *serverutil.NetAddress, kind string, list [][]byte) ([]byte, error) {
	s.Log(true, "Inv kind:"+kind+" sent to:", addrTo.String())
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
		na := serverutil.NewNetAddressIPPort(util.StringToNetIpAndPort(addr))
		_, err := s.sendInv(na, kind, list)
		if err == nil {
			nbRequestSucceeded++
		}
	}
	return float64(len(peers)) / float64(nbRequestSucceeded) * 100
}

//Receptionne une liste de hash de data (block || tx)
func (s *Server) handleInv(request []byte) {
	var payload serverutil.MsgInv
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}

	addr := payload.AddrSender.String()
	//augmente le nombre bytes recu depuis ce noeud.
	p, _ := s.GetPeer(addr)
	p.IncreaseBytesReceived(uint64(len(request)))
	s.AddPeer(p)
	s.Log(true, "Inv kind:"+payload.Kind+" received from :", addr)

	//si la requete inv envoie des hash de blocks
	if payload.Kind == "block" {
		var gbh *peerhistory.GetBlocksHistory
		var heightExpectedOfFirstElem = -1

		//on récupère un historique de demande de block en vers le noeud émetteur
		//de la requête inv courante.
		gbh = s.HistoryManager.GetBlock[addr].Select(true).SortByDate(true).First()
		if gbh != nil {
			//si l'intervalle demandé précédemment est supérieur ou égale au nombre maximal de block
			//qu'un noeud peut demander par requête
			if gbh.Message.Range[1]-gbh.Message.Range[0]+1 >= conf.MaxBlockPerMsg {
				//la hauteur de block attendu du premier element de la liste est
				//la hauteur initial demandé dans une precedente requête
				heightExpectedOfFirstElem = gbh.Message.Range[0]
			}
		}
		s.rangeBlockList(payload.AddrSender, payload.List, p, heightExpectedOfFirstElem)
	} else {
		s.rangeTxList(payload.AddrSender, payload.List)
	}
}

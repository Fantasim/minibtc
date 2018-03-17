package server

import (
	"log"
	conf "tway/config"
	"tway/serverutil"
	"tway/twayutil"
	"tway/util"
)

func (s *Server) askNewBlock(p *serverPeer, lastblock int) {
	var heightTo int
	if p == nil {
		listPeers := s.SelectPerfectPeers("getblock")
		if len(listPeers) == 0 {
			return
		}
		p = listPeers[0]
		lastblock = int(p.GetLastBlock())
	}
	if p != nil {
		//on récupère la plus haute hauteur de block demandé au réseau
		betterHeightAsked := s.HistoryManager.GetBetterHeightAsked()
		//si le noeud n'a pas demandé de nouveaux blocks ou
		//que la taille de la chain local est supérieur
		//à la dernière hauteur de block demandé.
		if betterHeightAsked == 0 || betterHeightAsked < s.chain.Height {
			betterHeightAsked = s.chain.Height
		}
		//si la hauteur de block du noeud emetteur de la version
		//est supérieur à la plus haute version demandé ou à la hauteur de la chain courante
		if betterHeightAsked < lastblock {
			//on récupère l'historique de demande de block ayant un intervalle [x; x + conf.MaxBlockPerMsg],
			//dont la valeur extérieur est la plus grande
			higherRangeAsked := s.HistoryManager.GetgetBlocksHistorysAskedWithMaxRange().HigherRange()
			var rangeExter int
			if higherRangeAsked != nil {
				rangeExter = higherRangeAsked.Message.Range[1]
			}
			//si la hauteur de block demandé ayant un intervalle maximum possible est inférieur
			//a la hauteur de la chaine courante
			if rangeExter <= s.chain.Height {
				heightTo = lastblock
				if heightTo-(betterHeightAsked) > conf.MaxBlockPerMsg {
					heightTo = conf.MaxBlockPerMsg + betterHeightAsked
				}
				na := serverutil.NewNetAddressIPPort(util.StringToNetIpAndPort(p.GetAddr()))
				_, err := s.sendAskBlocks(na, [2]int{s.chain.Height + 1, heightTo})
				//si la hauteur maximal demandé - la hauteur de chaine actuelle
				//est egale au nombre maximum de block demandable par requete
				//on set newFetchAtHeight à heightTo.
				/*
					ce qui fait que cette fonction sera appelée,
					lorsque la hauteur de la chaine sera égale à heightTo
					(fonction appelé dans ./block_manager.go -> BlockDownloaded)
				*/
				if err == nil && (heightTo-s.chain.Height) >= conf.MaxBlockPerMsg {
					s.newFetchAtHeight = s.chain.Height + conf.MaxBlockPerMsg
				}
			}
		}
	}
	return
}

func (s *Server) NewMsgAskBlock(rng [2]int) *serverutil.MsgAskBlocks {
	return &serverutil.MsgAskBlocks{s.ipStatus, rng}
}

func (s *Server) sendAskBlocks(addrTo *serverutil.NetAddress, rng [2]int) ([]byte, error) {
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
func (s *Server) handleAskBlocks(request []byte) {
	var payload serverutil.MsgAskBlocks
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}
	addr := payload.Addr.String()
	p, _ := s.GetPeer(addr)
	p.IncreaseBytesReceived(uint64(len(request)))
	s.AddPeer(p)
	s.HistoryManager.NewGetBlocksHistory(&payload, false)

	//récupère une liste de block dans un intervalle donné
	listBlock := s.chain.GetNBlocksNextToHeight(payload.Range[0], payload.Range[1]-payload.Range[0]+1)
	s.Log(true, "GetBlocks received from :", addr)
	s.Log(false, "height:", payload.Range[0]-1)
	s.Log(false, len(listBlock), "blocks found")

	//recupère une liste de hash depuis une list de block
	listHash := twayutil.GetListBlocksHashFromMap(listBlock)
	s.sendInv(payload.Addr, "block", listHash)
}

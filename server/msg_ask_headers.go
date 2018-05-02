package server

import (
	"bytes"
	"log"
	conf "tway/config"
	"tway/util"

	"tway/serverutil"
	"tway/twayutil"
)

func (s *Server) NewMsgAskHeaders(addrTo *serverutil.NetAddress, headHash []byte, stoppingHash []byte, count uint16) *serverutil.MsgAskHeaders {
	if len(headHash) == 0 {
		headHash = conf.GENESIS_BLOCK_PREVHASH
	}
	if len(stoppingHash) == 0 {
		stoppingHash = conf.GENESIS_BLOCK_PREVHASH
	}
	return &serverutil.MsgAskHeaders{s.ipStatus, addrTo, s.version, headHash, stoppingHash, count}
}

func (s *Server) sendAskHeaders(addrTo *serverutil.NetAddress, headHash []byte, stoppingHash []byte, count uint16) ([]byte, error) {
	s.Log(true, "GetHeaders sent to:", addrTo.String())
	askHeaders := s.NewMsgAskHeaders(addrTo, headHash, stoppingHash, count)
	//assigne en []byte la structure getblocks
	payload := gobEncode(*askHeaders)
	//on append la commande et le payload
	request := append(commandToBytes("getheaders"), payload...)
	err := s.sendData(addrTo.String(), request)

	if err != nil {
		go s.HistoryManager.NewGetHeadersHistory(askHeaders, false)
	}

	return request, err
}

//Retourne le pourcentage de succès des envoie de requêtes
func (s *Server) BootstrapGetHeaders(peers map[string]*serverPeer, headHash []byte, stoppingHash []byte, count uint16) float64 {
	lengthPeers := float64(len(peers))
	var nbRequestSucceeded = 0

	for addr, _ := range peers {
		na := serverutil.NewNetAddressIPPort(util.StringToNetIpAndPort(addr))
		_, err := s.sendAskHeaders(na, headHash, stoppingHash, count)
		if err == nil {
			nbRequestSucceeded++
		}
	}
	if lengthPeers == 0 {
		lengthPeers = 1
	}

	return lengthPeers / float64(nbRequestSucceeded) * 100
}

func (s *Server) filterHeadersListFromHeadHash(headBlock *twayutil.Block, count int, heightStart int) []serverutil.Header {
	if heightStart < count {
		count = heightStart
	}
	var headerList []serverutil.Header
	headerList = make([]serverutil.Header, count)
	b := headBlock
	var i int
	for i = count; i > 0; i-- {
		b = s.chain.GetBlockByHeight(heightStart - count + int(i))
		headerList[i-1] = serverutil.Header{heightStart - count + int(i), b.GetHash(), b.Header}
	}
	return headerList
}

func (s *Server) filterHeadersListFromStoppingHash(headBlock *twayutil.Block, count int, heightStart int) []serverutil.Header {
	if s.chain.Height < heightStart+count {
		count = s.chain.Height - heightStart + 1
	}

	var headerList []serverutil.Header
	headerList = make([]serverutil.Header, count)
	b := headBlock
	var i int
	for i = 0; i < count; i++ {
		b = s.chain.GetBlockByHeight(heightStart + i)
		headerList[i] = serverutil.Header{heightStart + i, b.GetHash(), b.Header}
	}
	return headerList
}

func (s *Server) GetBlockHeadersList(payload serverutil.MsgAskHeaders) []serverutil.Header {
	if payload.Count > conf.MaxHeadersPerMsg {
		payload.Count = conf.MaxHeadersPerMsg
	}

	var headBlock *twayutil.Block
	var headBlockHeight int
	var fromStopping bool
	var headerList []serverutil.Header

	//si il n'y a pas de headBlock demandé
	if bytes.Compare(payload.HeadHash, conf.GENESIS_BLOCK_PREVHASH) == 0 {
		//le headBlock n'étant pas defini.
		//le headblock est donc le head block de la chain
		headBlock = s.chain.GetLastBlock()
		headBlockHeight = s.chain.Height
		if bytes.Compare(payload.StoppingHash, conf.GENESIS_BLOCK_PREVHASH) != 0 {
			fromStopping = true
		}
	} else {
		headBlock, headBlockHeight = s.chain.GetBlockByHash(payload.HeadHash)
	}

	//si il n'y a pas de stopping hash defini
	if bytes.Compare(payload.StoppingHash, conf.GENESIS_BLOCK_PREVHASH) == 0 {
		if headBlock == nil {
			return headerList
		}
		headerList = s.filterHeadersListFromHeadHash(headBlock, int(payload.Count), headBlockHeight)
	} else {
		//on recupere le block et la hauteur correspondant au stopping block
		stoppingBlock, stoppingHeight := s.chain.GetBlockByHash(payload.StoppingHash)
		//si le stoppingblock existe
		if stoppingBlock != nil {
			//le nombre de header dans la liste correspond à :
			//hauteur du headHash - hauteur du stoppingHash
			var count int
			if fromStopping == true {
				count = int(payload.Count)
			} else {
				count = headBlockHeight - stoppingHeight + 1
			}
			if count > conf.MaxHeadersPerMsg {
				count = conf.MaxHeadersPerMsg
			}
			if fromStopping == true {
				headerList = s.filterHeadersListFromStoppingHash(stoppingBlock, count, stoppingHeight)
			} else {
				headerList = s.filterHeadersListFromHeadHash(headBlock, count, headBlockHeight)
			}
		}
	}
	return headerList
}

//Receptionne une demande de liste header de block
//voir structure MsgAskHeaders
func (s *Server) handleAskHeaders(request []byte) {
	var payload serverutil.MsgAskHeaders

	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}

	go s.HistoryManager.NewGetHeadersHistory(&payload, false)

	addr := payload.AddrSender.String()
	s.Log(true, "GetHeaders received from:", addr)

	p, _ := s.GetPeer(addr)
	p.IncreaseBytesReceived(uint64(len(request)))

	s.AddPeer(p)

	headerSlice := s.GetBlockHeadersList(payload)
	s.sendHeaders(payload.AddrReceiver, headerSlice, &payload)
}

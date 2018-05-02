package server

import (
	"bytes"
	"fmt"
	"time"
	conf "tway/config"
	"tway/serverutil"
	"tway/util"
)

//Recupère une liste d'address non traité
//Ce sont des pairs avec qui le noeud courant n'a ni vérifié l'existence avec un ping,
//ni échangé les versions.
func (s *Server) ListOfUntreatedPeers() listOfPeers {
	ret := make(listOfPeers)
	s.peers.Range(func(key, val interface{}) bool {
		p := val.(*serverPeer)
		if p.GetLastPingSentTime() == 0 && p.IsVersionSent() == false {
			ret[p.GetAddr()] = p
		}
		return true
	})
	return ret
}

//Recupère une liste d'address dite "de confiance"
//ce sont des pairs avec qui le noeud courant a échangé sa version
func (s *Server) ListOfTrustedPeers() listOfPeers {
	ret := make(listOfPeers)
	s.peers.Range(func(key, val interface{}) bool {
		p := val.(*serverPeer)
		if p.IsVersionSent() == true && p.IsVerAckReceived() == true {
			ret[p.GetAddr()] = p
		}
		return true
	})

	return ret
}

func (s *Server) GetCloserAndSafestPeers() listOfPeers {
	ret := make(listOfPeers)

	s.peers.Range(func(key, val interface{}) bool {
		p := val.(*serverPeer)
		if p.IsVersionSent() == true && p.IsVerAckReceived() == true {
			ret[p.GetAddr()] = p
		}
		return true
	})

	return ret
}

func (s *Server) GetListOfTrustedMainNode() listOfPeers {
	ret := make(listOfPeers)

	s.peers.Range(func(key, val interface{}) bool {
		p := val.(*serverPeer)
		addr := p.GetAddr()
		na := serverutil.NewNetAddressIPPort(util.StringToNetIpAndPort(addr))
		na.IsEqual(GetMainNode())
		if na.IsEqual(GetMainNode()) == true && p.IsVersionSent() == true && p.IsVerAckReceived() == true {
			ret[addr] = p
		}
		return true
	})

	return ret
}

//Cette fonction est appelé dès lors que le noeud recoit une liste de nouvelles addresses
//Pour chaque adresse recu, elle va lui envoyer un ping puis attendre un certain délais
//pour la reception d'un pong provenant de ce même noeud.
//Si elle recoit un pong dans les conf.TimeInSecondAfterPingWithoutPongToRemove seconds
//Elle lui envoie une version du noeud.
func (s *Server) treatPeersAfterPong(unTreatedPeers map[string]*serverPeer) {
	//pour chaque pair non traité
	for newPeerAddr, _ := range unTreatedPeers {
		na := serverutil.NewNetAddressIPPort(util.StringToNetIpAndPort(newPeerAddr))
		//on envoie un ping
		s.sendPing(na)
	}

	//tant que la liste contient des pairs non traités, on boucle a l'infini
	for len(unTreatedPeers) > 0 {
		//pour chaque pair non traité
		for unTreatedaddr, p := range unTreatedPeers {
			nowUnixNano := time.Now().UnixNano()
			//si le ping n'obtient pas de pong réponse dans les 5 secondes suivant l'envoie du ping
			condition1 := ((nowUnixNano - p.GetLastPingSentTime()) / int64(time.Second)) >= conf.TimeInSecondAfterPingWithoutPongToRemove
			condition2 := ((p.GetLastPongReceivedTime() - p.GetLastPingSentTime()) / int64(time.Second)) >= conf.TimeInSecondAfterPingWithoutPongToRemove
			//si le ping obtient une réponse pong dans les 5 secondes suivant l'envoie du ping
			condition3 := p.GetLastPongReceivedTime() > p.GetLastPingSentTime() && ((p.GetLastPongReceivedTime()-p.GetLastPingSentTime())/int64(time.Second)) < conf.TimeInSecondAfterPingWithoutPongToRemove
			if condition1 || condition2 {
				s.Log(true, "pong not received from", unTreatedaddr)
				delete(unTreatedPeers, unTreatedaddr)
			} else if condition3 {
				//si l'on a pas déjà envoyé notre version au pair
				p, _ := s.GetPeer(unTreatedaddr)
				if p.IsVersionSent() == false {
					//on envoie notre version au pair
					s.sendVersion(serverutil.NewNetAddressIPPort(util.StringToNetIpAndPort(unTreatedaddr)))
				}
				//on supprime le pair de la liste des pairs non traités
				delete(unTreatedPeers, unTreatedaddr)
			}
		}
		//on sleep durant 100 MS
		time.Sleep(100 * time.Millisecond)
	}
}

/*IMPROVE_LATER*/
//kind :
//getblock
func (s *Server) SelectPerfectPeers(kind string) sliceOfPeers {
	switch kind {
	case "getblock":
		peers := s.syncMapToListOfPeers(s.peers).GetPeersBasedOnHeight(s.chain.Height + 1).ListOfPeersToSlice()
		peers.SortByBytesReceived(false)
		return peers
	default:
		return nil
	}
	return nil
}

func (s *Server) SelectPerfectPeersHavingABlock(hash []byte) sliceOfPeers {
	var ret sliceOfPeers

	headersHistory := s.HistoryManager.Headers
HEADERHISTORY:
	for addr, listHH := range headersHistory {

		for _, hh := range listHH {

			for _, h := range hh.Message.List {
				if bytes.Compare(h.Hash, hash) == 0 || bytes.Compare(h.Header.HashPrevBlock, hash) == 0 {
					na, err := serverutil.NewNetAddressByString(addr)
					if err == nil {
						ret = append(ret, NewServerPeer(na))
					} else {
						fmt.Println("ERROR IN SelectPerfectPeersHavingABlock")
					}
					continue HEADERHISTORY
				}
			}
		}
	}
	return ret
}

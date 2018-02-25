package server

import (
	"tway/util"
	"time"
	conf "tway/config"
)

//Recupère une liste d'address non traité
//Ce sont des pairs avec qui le noeud courant n'a ni vérifié l'existence avec un ping,
//ni échangé les versions.
func (s *Server) ListOfUntreatedPeers() map[string]*serverPeer {
	ret := make(map[string]*serverPeer)

	for addr, p := range s.peers {
		if p.GetLastPingSentTime() == 0 && p.IsVersionSent() == false {
			ret[addr] = p
		}
	}
	return ret
}

//Recupère une liste d'address dite "de confiance"
//ce sont des pairs avec qui le noeud courant a échangé sa version
func (s *Server) ListOfTrustedPeers() map[string]*serverPeer {
	ret := make(map[string]*serverPeer)

	for addr, p := range s.peers {
		if p.IsVersionSent() == true && p.IsVerAckReceived() == true {
			ret[addr] = p
		}
	}
	return ret
}

//Cette fonction est appelé dès lors que le noeud recoit une liste de nouvelles addresses
//Pour chaque adresse recu, elle va lui envoyer un ping puis attendre un certain délais 
//pour la reception d'un pong provenant de ce même noeud.
//Si elle recoit un pong dans les conf.TimeInSecondAfterPingWithoutPongToRemove seconds
//Elle lui envoie une version du noeud.
func (s *Server) treatPeersAfterPong(unTreatedPeers map[string]*serverPeer){
	defer s.addrMu.Unlock()
	//pour chaque pair non traité
	for newPeerAddr, _ := range unTreatedPeers {
		na := NewNetAddressIPPort(util.StringToNetIpAndPort(newPeerAddr))
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
			condition3 := p.GetLastPongReceivedTime() > p.GetLastPingSentTime() && ((p.GetLastPongReceivedTime() - p.GetLastPingSentTime()) / int64(time.Second)) < conf.TimeInSecondAfterPingWithoutPongToRemove
			if condition1 || condition2 {
				s.Log(true, "pong not received from", unTreatedaddr)
				delete(unTreatedPeers, unTreatedaddr)
			} else if condition3 {
				//si l'on a pas déjà envoyé notre version au pair
				if s.peers[unTreatedaddr].IsVersionSent() == false {
					//on envoie notre version au pair 
					s.sendVersion(NewNetAddressIPPort(util.StringToNetIpAndPort(unTreatedaddr)))
				}
				//on supprime le pair de la liste des pairs non traités
				delete(unTreatedPeers, unTreatedaddr)
			}
		}
		//on sleep durant 500 MS
		time.Sleep(500 * time.Millisecond)
	}
}

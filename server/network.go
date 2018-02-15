package server

import (
	"tway/util"
	"time"
	conf "tway/config"
)

func (s *Server) ListOfUntreatedPeers() map[string]*serverPeer {
	ret := make(map[string]*serverPeer)

	for addr, p := range s.peers {
		if p.GetLastPingSentTime() == 0 && p.IsVersionSent() == false {
			ret[addr] = p
		}
	}
	return ret
}

func (s *Server) ListOfTrustedPeers() map[string]*serverPeer {
	ret := make(map[string]*serverPeer)

	for addr, p := range s.peers {
		if p.IsVersionSent() == true && p.IsVerAckReceived() == true {
			ret[addr] = p
		}
	}
	return ret
}

func (s *Server) treatPeersAfterPong(unTreatedPeers map[string]*serverPeer){
	for newPeerAddr, _ := range unTreatedPeers {
		na := NewNetAddressIPPort(util.StringToNetIpAndPort(newPeerAddr))
		s.sendPing(na)
	}
	
	for len(unTreatedPeers) > 0 {
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
				if s.peers[unTreatedaddr].IsVersionSent() == false {
					s.sendVersion(NewNetAddressIPPort(util.StringToNetIpAndPort(unTreatedaddr)))
				}
				delete(unTreatedPeers, unTreatedaddr)
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
}
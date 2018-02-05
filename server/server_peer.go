package server

import (
	"tway/peer"
	"sync"
)

type serverPeer struct {
	*peer.Peer

	mu 		sync.Mutex
	mainNode bool
}

func NewServerPeer(addr string) *serverPeer {
	return &serverPeer{
		Peer: peer.NewPeer(addr),
		mainNode: GetMainNode().String() == addr, 
	}
}

func (sp *serverPeer) IsEqual(sp2 *serverPeer) bool {
	return sp.Peer.GetAddr() == sp2.Peer.GetAddr() 
}

func (sp *serverPeer) IsMainNode() bool {
	return sp.mainNode
}

func (sp *serverPeer) VerAckReceived(){
	sp.mu.Lock()
	sp.Peer.VerAckReceived()
	defer sp.mu.Unlock()
}
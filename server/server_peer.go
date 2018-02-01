package server

import (
	"tway/peer"
)

type serverPeer struct {
	*peer.Peer
	mainNode bool
}

func NewServerPeer(addr string) *serverPeer {
	return &serverPeer{peer.NewPeer(addr), GetMainNode().String() == addr}
}

func (sp *serverPeer) IsEqual(sp2 *serverPeer) bool {
	return sp.Peer.GetAddr() == sp2.Peer.GetAddr() 
}

func (sp *serverPeer) IsMainNode() bool {
	return sp.mainNode
}
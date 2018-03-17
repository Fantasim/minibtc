package server

import (
	"sync"
	"tway/peer"

	"github.com/bradfitz/slice"
)

type listOfPeers map[string]*serverPeer
type sliceOfPeers []*serverPeer

type serverPeer struct {
	*peer.Peer

	mu       sync.Mutex
	mainNode bool
}

func NewServerPeer(addr string) *serverPeer {
	return &serverPeer{
		Peer:     peer.NewPeer(addr),
		mainNode: GetMainNode().String() == addr,
	}
}

func (sp *serverPeer) IsEqual(sp2 *serverPeer) bool {
	return sp.Peer.GetAddr() == sp2.Peer.GetAddr()
}

func (sp *serverPeer) IsMainNode() bool {
	return sp.mainNode
}

func (sp *serverPeer) VerAckReceived() {
	sp.mu.Lock()
	sp.Peer.VerAckReceived()
	defer sp.mu.Unlock()
}

func (s *Server) syncMapToListOfPeers(m sync.Map) listOfPeers {
	ret := make(listOfPeers, 0)
	m.Range(func(key, val interface{}) bool {
		p := val.(*serverPeer)
		ret[key.(string)] = p
		return true
	})
	return ret
}

/*IMPROVE_LATER*/
//récupère une liste de pair ayant une hauteur de block supérieur à minHeight
func (list listOfPeers) GetPeersBasedOnHeight(minHeight int) listOfPeers {
	ret := make(listOfPeers, 0)

	for addr, p := range list {
		if p.GetLastBlock() >= int64(minHeight) {
			ret[addr] = p
		}
	}
	return ret
}

func (list listOfPeers) ListOfPeersToSlice() sliceOfPeers {
	ret := make([]*serverPeer, 0)

	for _, item := range list {
		ret = append(ret, item)
	}
	return ret
}

func (slc sliceOfPeers) SliceOfPeersToMap() listOfPeers {
	ret := make(map[string]*serverPeer)

	for _, item := range slc {
		ret[item.GetAddr()] = item
	}
	return ret
}

//Trie un slice de pair selon le nombre de bytes que le pair
//courant a reçu de leur part.
func (slc sliceOfPeers) SortByBytesReceived(desc bool) {
	//on les tries par date de creation
	slice.Sort(slc[:], func(i, j int) bool {
		if desc == false {
			return slc[i].GetBytesReceived() > slc[j].GetBytesReceived()
		} else {
			return slc[i].GetBytesReceived() < slc[j].GetBytesReceived()
		}
	})
}

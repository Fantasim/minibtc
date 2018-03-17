package peerhistory

import (
	"time"
	conf "tway/config"

	"github.com/bradfitz/slice"

	"tway/serverutil"
)

type ListBlocksHistory []GetBlocksHistory

//structure representant un historique de getblocks
type GetBlocksHistory struct {
	Message *serverutil.MsgAskBlocks
	Date    time.Time
	Sent    bool
}

//Récupère la plus grande hauteur de block demandé
func (hm *HistoryManager) GetBetterHeightAsked() int {
	var betterHeight int = 0
	for _, gb := range hm.GetBlock {
		for _, h := range gb {
			if betterHeight < h.Message.Range[1] {
				betterHeight = h.Message.Range[1]
			}
		}
	}
	return betterHeight
}

//Recupère une liste de requete getblocks ayant un intervalle de block
//maximum possible soit Range[1] - Range[0] == conf.MaxBlockPerMsg
//le nombre maximale possible de block demandable par requete getblocks
func (hm *HistoryManager) GetgetBlocksHistorysAskedWithMaxRange() ListBlocksHistory {
	var list ListBlocksHistory
	for _, gb := range hm.GetBlock {
		for _, h := range gb {
			if h.ContainAFullRange() {
				list = append(list, h)
			}
		}
	}
	return list.Select(true)
}

//Créer un historique de requete getblocks
func (hm *HistoryManager) NewGetBlocksHistory(msg *serverutil.MsgAskBlocks, sent bool) {
	hm.muGetBlock.Lock()
	defer hm.muGetBlock.Unlock()
	addr := msg.Addr.String()
	hm.GetBlock[addr] = append(hm.GetBlock[addr], GetBlocksHistory{msg, time.Now(), sent})
}

//Fonction retournant une liste d'historique de requete getblocks
//ayant été envoyés par le noeud courant ou reçus.
//sent == true | envoyées
//sent == false | reçus
func (list ListBlocksHistory) Select(sent bool) ListBlocksHistory {
	var ret ListBlocksHistory
	for _, item := range list {
		if item.Sent == sent {
			ret = append(ret, item)
		}
	}
	return ret
}

//trie par date de creation une liste d'historique de requete getblocks
func (list ListBlocksHistory) SortByDate(desc bool) ListBlocksHistory {
	ret := make(ListBlocksHistory, len(list))
	copy(ret, list)

	//on les tries par date de creation
	slice.Sort(ret[:], func(i, j int) bool {
		if desc == false {
			return ret[i].Date.After(ret[j].Date)
		} else {
			return ret[i].Date.Before(ret[j].Date)
		}
	})
	return ret
}

func (list ListBlocksHistory) First() *GetBlocksHistory {
	if len(list) > 0 {
		return &list[0]
	}
	return nil
}

//recupere l'historique de requete getblocks dans une liste ayant :
// - la hauteur demandé la plus haute
// - un intervalle de block maximum possible soit
// Range[1] - Range[0] == conf.MaxBlockPerMsg
func (list ListBlocksHistory) HigherRange() *GetBlocksHistory {
	var higher = 0
	var ret *GetBlocksHistory

	if len(list) == 0 {
		return nil
	}

	for _, item := range list {
		if item.Message.Range[1] > higher && item.ContainAFullRange() == true {
			ret = &item
			higher = item.Message.Range[1]
		}
	}
	return ret
}

//Retourne true si un historique de requete getblocks
//contient un intervalle de block maximum possible
//soit Range[1] - Range[0] + 1 == conf.MaxBlockPerMsg
func (gbh *GetBlocksHistory) ContainAFullRange() bool {
	if gbh == nil {
		return false
	}
	return gbh.Message.Range[1]-gbh.Message.Range[0]+1 == conf.MaxBlockPerMsg
}

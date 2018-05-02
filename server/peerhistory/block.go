package peerhistory

import (
	"fmt"
	"time"
	conf "tway/config"

	"github.com/bradfitz/slice"

	"tway/serverutil"
)

type GetBlock map[string]ListBlocksHistory
type ListBlocksHistory []GetBlocksHistory

//structure representant un historique de getblocks
type GetBlocksHistory struct {
	Date    time.Time
	Message *serverutil.MsgAskBlocks
	Sent    bool
}

//Créer un historique de requete getblocks
func (hm *HistoryManager) NewGetBlocksHistory(msg *serverutil.MsgAskBlocks, sent bool) {
	hm.muGetBlock.Lock()
	defer hm.muGetBlock.Unlock()
	addr := msg.Addr.String()
	gbh := GetBlocksHistory{time.Now(), msg, sent}
	hm.GetBlock[addr] = append(hm.GetBlock[addr], gbh)
	hm.Log("getblocks request", gbh.String())
}

//Récupère la plus grande hauteur de block demandé
func (gb GetBlock) GetBetterHeightAsked() int {
	var betterHeight int = 0
	for _, gbh := range gb {
		for _, h := range gbh {
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
func (gb GetBlock) GetgetBlocksHistorysAskedWithMaxRange() ListBlocksHistory {
	var list ListBlocksHistory
	for _, gbh := range gb {
		for _, h := range gbh {
			if h.ContainAFullRange() {
				list = append(list, h)
			}
		}
	}
	return list.Select(true)
}

func (gb GetBlock) GetBlocksHistoryByAddr(addr string) ListBlocksHistory {
	d, exist := gb[addr]
	if exist == false {
		var ret ListBlocksHistory
		return ret
	}
	return d
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
		if desc == true {
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

func (gbh *GetBlocksHistory) String() string {
	return fmt.Sprintf("{Date: %s, Message: %s, Sent: %t}", gbh.Date.Format("15:04:05.000000"), gbh.Message.String(), gbh.Sent)
}

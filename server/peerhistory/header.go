package peerhistory

import (
	"bytes"
	"fmt"
	"time"
	"tway/serverutil"

	"github.com/bradfitz/slice"
)

type Headers map[string]ListHeadersHistory
type ListHeadersHistory []HeadersHistory

type HeadersHistory struct {
	Date    time.Time
	Message *serverutil.MsgHeaders
}

//Créer un historique de requete headers
func (hm *HistoryManager) NewHeadersHistory(msg *serverutil.MsgHeaders) {
	hm.muGetHeader.Lock()
	defer hm.muGetHeader.Unlock()
	addr := msg.AddrSender.String()
	hh := HeadersHistory{time.Now(), msg}
	hm.Headers[addr] = append(hm.Headers[addr], hh)
	hm.Log("headers request:", hh.String())
}

func (h Headers) GetListHeadersHistoryByAddr(addr string) ListHeadersHistory {
	d, exist := h[addr]
	if exist == false {
		var ret ListHeadersHistory
		return ret
	}
	return d
}

//Cette requête retourne le nombre de pair ayant répondu à une requete getheaders.
func (h Headers) CountNbUniquePeerHavingAnsweredFromAGetHeadersRequest(stoppingHash, headHash []byte, count uint16) int {
	var ret int
	for _, list := range h {
		for _, headerReq := range list {
			if bytes.Compare(headerReq.Message.GetHeadersOrigin.HeadHash, headHash) == 0 &&
				bytes.Compare(headerReq.Message.GetHeadersOrigin.StoppingHash, stoppingHash) == 0 &&
				headerReq.Message.GetHeadersOrigin.Count == count {
				ret++
				break
			}
		}
	}
	return ret
}

//Fonction retournant une liste d'historique de requete headers
//dont le nombre de headers dans la liste est inférieur ou egale a count
func (list ListHeadersHistory) SelectByCount(count int) ListHeadersHistory {
	var ret ListHeadersHistory

	for _, item := range list {
		if len(item.Message.List) <= count {
			ret = append(ret, item)
		}
	}
	return ret
}

//Fonction retournant une liste d'historique de requete headers
//dont le nombre de headers dans la liste est inférieur ou egale a count
func (list ListHeadersHistory) GetLonguestListHeaderLength() int {
	var longuest = 0
	for _, item := range list {
		if len(item.Message.List) > longuest {
			longuest = len(item.Message.List)
		}
	}
	return longuest
}

//trie par date de creation une liste d'historique de requete getheaders
func (list ListHeadersHistory) SortByDate(desc bool) ListHeadersHistory {
	ret := make(ListHeadersHistory, len(list))
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

func (list ListHeadersHistory) First() *HeadersHistory {
	if len(list) > 0 {
		return &list[0]
	}
	return nil
}

//Fonction retournant une liste d'historique de requete headers
//comportant le stoppingHash passé en parametre
func (list ListHeadersHistory) SelectByStoppingHash(stoppingHash []byte) ListHeadersHistory {
	var ret ListHeadersHistory
	if len(list) == 0 {
		return ret
	}
	for _, item := range list {
		if bytes.Compare(item.Message.List[0].Hash, stoppingHash) == 0 {
			ret = append(ret, item)
		}
	}
	return ret
}

//Fonction retournant une liste d'historique de requete headers
//comportant le headHash passé en parametre
func (list ListHeadersHistory) SelectByHeadHash(headHash []byte) ListHeadersHistory {
	var ret ListHeadersHistory
	if len(list) == 0 {
		return ret
	}
	for _, item := range list {
		if bytes.Compare(item.Message.List[len(item.Message.List)-1].Hash, headHash) == 0 {
			ret = append(ret, item)
		}
	}
	return ret
}

func (hh *HeadersHistory) String() string {
	return fmt.Sprintf("{Date: %s, Message: %s}", hh.Date.Format("15:04:05.000000"), hh.Message.String())
}

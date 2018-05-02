package peerhistory

import (
	"bytes"
	"fmt"
	"time"
	"tway/serverutil"

	"github.com/bradfitz/slice"
)

type GetHeader map[string]ListGetHeadersHistory
type ListGetHeadersHistory []GetHeadersHistory

type GetHeadersHistory struct {
	Message *serverutil.MsgAskHeaders
	Date    time.Time
	Sent    bool
	ID      int
}

//Créer un historique de requete getheaders
func (hm *HistoryManager) NewGetHeadersHistory(msg *serverutil.MsgAskHeaders, sent bool) {
	hm.muGetHeader.Lock()
	defer hm.muGetHeader.Unlock()
	var addr string
	if sent == true {
		addr = msg.AddrReceiver.String()
	} else {
		addr = msg.AddrSender.String()
	}
	id := len(hm.GetHeader[addr])
	ghh := GetHeadersHistory{msg, time.Now(), sent, id}
	hm.GetHeader[addr] = append(hm.GetHeader[addr], ghh)
	hm.Log("getheaders request", ghh.String())
}

func (gh GetHeader) GetListGetHeadersHistoryByAddr(addr string) ListGetHeadersHistory {
	d, exist := gh[addr]
	if exist == false {
		var ret ListGetHeadersHistory
		return ret
	}
	return d
}

//Fonction retournant une liste d'historique de requete getheaders
//ayant été envoyés par le noeud courant ou reçus.
//sent == true | envoyées
//sent == false | reçus
func (list ListGetHeadersHistory) SelectBySent(sent bool) ListGetHeadersHistory {
	var ret ListGetHeadersHistory

	for _, item := range list {
		if item.Sent == sent {
			ret = append(ret, item)
		}
	}
	return ret
}

//Fonction retournant une liste d'historique de requete getheaders
//comportant le stoppingHash passé en parametre
func (list ListGetHeadersHistory) SelectByStoppingHash(stoppingHash []byte) ListGetHeadersHistory {
	var ret ListGetHeadersHistory
	for _, item := range list {
		if bytes.Compare(item.Message.StoppingHash, stoppingHash) == 0 {
			ret = append(ret, item)
		}
	}
	return ret
}

//Fonction retournant une liste d'historique de requete getheaders
//comportant le stoppingHash passé en parametre
func (list ListGetHeadersHistory) SelectByHeadHash(headHash []byte) ListGetHeadersHistory {
	var ret ListGetHeadersHistory
	for _, item := range list {
		if bytes.Compare(item.Message.HeadHash, headHash) == 0 {
			ret = append(ret, item)
		}
	}
	return ret
}

//Fonction retournant une liste d'historique de requete getheaders
//comportant a un count supérieur egale a la valeur passé en param
func (list ListGetHeadersHistory) SelectByCount(count int) ListGetHeadersHistory {
	var ret ListGetHeadersHistory
	for _, item := range list {
		if item.Message.Count == uint16(count) {
			ret = append(ret, item)
		}
	}
	return ret
}

//trie par date de creation une liste d'historique de requete getheaders
func (list ListGetHeadersHistory) SortByDate(desc bool) ListGetHeadersHistory {
	ret := make(ListGetHeadersHistory, len(list))
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

func (list ListGetHeadersHistory) First() *GetHeadersHistory {
	if len(list) > 0 {
		return &list[0]
	}
	return nil
}

func (ghh *GetHeadersHistory) String() string {
	return fmt.Sprintf("{Date:%s, ID: %d, Message: %s, Sent: %t}", ghh.Date.Format("15:04:05.000000"), ghh.ID, ghh.Message.String(), ghh.Sent)
}

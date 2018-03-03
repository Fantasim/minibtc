package server

import (
	"time"
	"github.com/bradfitz/slice"
	conf "tway/config"
	"sync"
)

type listBlocksHistory []getBlocksHistory

type historyManager struct {
	GetBlock map[string]listBlocksHistory
	mu		sync.Mutex
}

type getBlocksHistory struct {
	Message	*MsgAskBlocks
	Date 	time.Time
	Sent	bool
}

func NewHistoryManager() *historyManager {
	return &historyManager{
		GetBlock: make(map[string]listBlocksHistory),
	}
}

func (hm *historyManager) GetBetterHeightAsked() int {
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

func (hm *historyManager) GetgetBlocksHistorysAskedWithMaxRange() listBlocksHistory {
	var list listBlocksHistory
	for _, gb := range hm.GetBlock {
		for _, h := range gb {
			if h.ContainAFullRange() {
				list = append(list, h)
			}
		}
	}
	return list.Select(true)
}

func (hm *historyManager) NewGetBlocksHistory(msg *MsgAskBlocks, sent bool){
	hm.mu.Lock()
	defer hm.mu.Unlock()
	addr := msg.Addr.String()
	hm.GetBlock[addr] = append(hm.GetBlock[addr], getBlocksHistory{msg, time.Now(), sent})
}

func (list listBlocksHistory) Select(sent bool) listBlocksHistory {
	var ret listBlocksHistory
	for _, item := range list {
		if item.Sent == sent {
			ret = append(ret, item)
		}
	}
	return ret
}

func (list listBlocksHistory) sortByDate(desc bool) listBlocksHistory {
	ret := make(listBlocksHistory, len(list))
	copy(ret, list)

	//on les tries par date de creation
	slice.Sort(ret[:], func(i,j int) bool {
		if desc == false {
			return ret[i].Date.After(ret[j].Date)
		} else {
			return ret[i].Date.Before(ret[j].Date)
		}
	})
	return ret
}

func (list listBlocksHistory) first() *getBlocksHistory {
	if len(list) > 0 {
		return &list[0]
	}
	return nil
}

func (list listBlocksHistory) higherRange() *getBlocksHistory {
	var higher = 0
	var ret *getBlocksHistory
	
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

func (gbh *getBlocksHistory) ContainAFullRange() bool {
	if gbh == nil {
		return false
	}
	return gbh.Message.Range[1] - gbh.Message.Range[0] + 1 == conf.MaxBlockPerMsg	
}
package server

import (
	"time"
)

type historyManager struct {
	GetBlock map[string][]getBlocksHistory
}

type getBlocksHistory struct {
	Message	*MsgAskBlocks
	Date 	time.Time
	Sent	bool
}

func NewHistoryManager() *historyManager {
	return &historyManager{
		GetBlock: make(map[string][]getBlocksHistory),
	}
}

func (hm *historyManager) NewGetBlocksHistory(msg *MsgAskBlocks, sent bool){
	addr := msg.Addr.String()
	hm.GetBlock[addr] = append(hm.GetBlock[addr], getBlocksHistory{msg, time.Now(), sent})
}

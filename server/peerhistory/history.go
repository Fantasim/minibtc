package peerhistory

import (
	"sync"
)

//structure representant le manager des historiques de requetes
type HistoryManager struct {
	GetBlock   map[string]ListBlocksHistory
	muGetBlock sync.Mutex

	GetHeader   map[string]ListHeadersHistory
	muGetHeader sync.Mutex
}

//Creer un nouveau manager d'historique
func NewHistoryManager() *HistoryManager {
	return &HistoryManager{
		GetBlock:  make(map[string]ListBlocksHistory),
		GetHeader: make(map[string]ListHeadersHistory),
	}
}

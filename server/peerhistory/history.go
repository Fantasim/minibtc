package peerhistory

import (
	"fmt"
	"sync"
)

//structure representant le manager des historiques de requetes
type HistoryManager struct {
	GetBlock   GetBlock
	muGetBlock sync.Mutex

	GetHeader   GetHeader
	muGetHeader sync.Mutex

	Headers Headers

	Version   Version
	muVersion sync.Mutex

	log bool
}

//Creer un nouveau manager d'historique
func NewHistoryManager(log bool) *HistoryManager {
	return &HistoryManager{
		GetBlock:  make(GetBlock),
		GetHeader: make(GetHeader),
		Version:   make(Version),
		Headers:   make(Headers),
		log:       log,
	}
}

func (hm *HistoryManager) Log(title, content string) {
	if hm.log == true {
		fmt.Printf("HISTORY %s: %s\n", title, content)
	}
}

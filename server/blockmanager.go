package server

import (
	"sync"
	"tway/twayutil"
	"encoding/hex"
	"time"
)

type DownloadInformations struct {
	sp 		*serverPeer
	start 	time.Time
	block		*twayutil.Block //==nil if block hasn't been still downloaded
}

type blockManager struct {
	download		map[string]*DownloadInformations
	mu 				sync.Mutex
}

func NewBlockManager() *blockManager {
	return &blockManager{
		download: make(map[string]*DownloadInformations),
	}
}

func (bm *blockManager) BlockDownloaded(new *twayutil.Block){
	hash := hex.EncodeToString(new.GetHash())
	bm.mu.Lock()
	defer bm.mu.Unlock()
	if bm.download[hash] == nil {
		return 
	} else {
		bm.download[hash].block = new
	}
}

func (bm *blockManager) StartDownloadBlock(hash string, sp *serverPeer){
	bm.mu.Lock()
	defer bm.mu.Unlock()
	bm.download[hash] = &DownloadInformations{sp, time.Now(), nil}
}
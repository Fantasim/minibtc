package server

import (
	"sync"
	"tway/twayutil"
	"encoding/hex"
	"time"
	b "tway/blockchain"
	"fmt"
	"bytes"
)

type DownloadInformations struct {
	sp 		*serverPeer
	start 	int64 //ns
	receivedAt int64 //ns
	block		*twayutil.Block //==nil if block hasn't been still downloaded
	timeToRetry	int64 //ns
}

type blockManager struct {
	download		map[string]*DownloadInformations
	mu 				sync.Mutex
	chain			*b.Blockchain
}

func NewBlockManager() *blockManager {
	return &blockManager{
		download: make(map[string]*DownloadInformations),
		chain: b.BC,
	}
}

func (bm *blockManager) IsDownloaded(hash string) bool {
	return bm.download[hash] != nil && bm.download[hash].block != nil
}

func (bm *blockManager) IsDownloading(hash string) bool {
	return bm.download[hash] != nil && bm.download[hash].block == nil
}

func (bm *blockManager) BlockDownloaded(new *twayutil.Block){
	hash := hex.EncodeToString(new.GetHash())	
	bm.mu.Lock()
	defer bm.mu.Unlock()
	if bm.download[hash] == nil {
		return
		//si le block recu est nil et que le telechargement du block est enregistré
	} else if bm.download[hash] != nil && new == nil {
		//si le block est en cours de téléchargement
		if bm.IsDownloading(hash) == true && bm.IsDownloaded(hash) == false {
			//on supprime le téléchargement
			delete(bm.download, hash)
		}
		//new block is nil
		return
	}

	if bm.download[hash].timeToRetry == 0 {
		err := bm.chain.CheckNewBlock(new)
		if err == nil || err.Error() == b.NO_NEXT_TO_TIP_ERROR {
			//
		} else {
			fmt.Println(err)
			delete(bm.download, hash)
			return
		}
	}

	lastChainBlock := b.BC.GetLastBlock()
	//SYSTEM ERROR - highly possible
	//if block is not next to last block added in chain
	if bytes.Compare(lastChainBlock.GetHash(), new.Header.HashPrevBlock) != 0 {
		averageTimeToDownloadBlock := GetAverageTimeToDownloadABlock(bm.download)
		if averageTimeToDownloadBlock == 0 {
			//500 MS
			averageTimeToDownloadBlock = 1000000 * 500
		}
		bm.download[hash].timeToRetry = time.Now().Add(time.Nanosecond * time.Duration(averageTimeToDownloadBlock)).UnixNano()
		go func(){
			time.Sleep(time.Nanosecond * time.Duration(averageTimeToDownloadBlock))
			bm.BlockDownloaded(new)
		}()
		return
	}
	err := bm.chain.AddBlock(new)
	if err == nil {
		fmt.Printf("block %d - %s successfully added on chain", bm.chain.Height, hex.EncodeToString(new.GetHash()))
		bm.download[hash].block = new
		bm.download[hash].receivedAt = time.Now().UnixNano()
	}
}

func (bm *blockManager) StartDownloadBlock(hash string, sp *serverPeer){
	bm.mu.Lock()
	defer bm.mu.Unlock()
	bm.download[hash] = &DownloadInformations{
		sp: sp, 
		start: time.Now().UnixNano(), 
	}
}
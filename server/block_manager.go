package server

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
	b "tway/blockchain"
	conf "tway/config"
	"tway/twayutil"
)

type listDownloadInformation []DownloadInformations

type DownloadInformations struct {
	sp               *serverPeer
	start            int64           //ns
	receivedAt       int64           //ns
	addedAt          int64           //ns
	expectedHeight   int64           //expected height for this block
	block            *twayutil.Block //==nil if block hasn't been still downloaded
	unConfirmedBlock *twayutil.Block
	timeToRetry      int64 //ns
	nbTry            int64
}

type blockManager struct {
	NewBlock chan *twayutil.Block
	download map[string]*DownloadInformations
	mu       sync.Mutex
	chain    *b.Blockchain
	log      bool
	orphanMu sync.Mutex
}

func NewBlockManager(log, mining bool) *blockManager {
	return &blockManager{
		NewBlock: make(chan *twayutil.Block),
		download: make(map[string]*DownloadInformations),
		chain:    *&b.BC,
		log:      log,
	}
}

//Cette fonction retourne true si le block lié au hash a été téléchargé
func (bm *blockManager) IsDownloaded(hash string) bool {
	return bm.download[hash] != nil && bm.download[hash].block != nil
}

//Cette fonction retourne true si le block lié au hash est en cours de téléchargement
func (bm *blockManager) IsDownloading(hash string) bool {
	return bm.download[hash] != nil && bm.download[hash].block == nil
}

//Cette fonction est appelé dans le fonction handle block [/block.go]
//Elle permet de receptionner un block téléchargé et de controler chaque partie du block
//permettant ainsi de le rejeter ou de l'ajouter a la blockchain locale.
//elle met egalement a jour le blockmanager en fonction du resultat
func (bm *blockManager) BlockDownloaded(new *twayutil.Block, s *Server) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	//si le block est vide
	if new == nil {
		bm.Log(false, "block is nil")
		return
	}
	//hash du block recu
	hash := hex.EncodeToString(new.GetHash())

	//si le block recu n'existe pas dans la liste des blocks en cours de téléchargement (si le block n'a pas été demandé)
	if bm.download[hash] == nil {
		bm.Log(false, "download information is not exist")
		return
	}

	if bm.download[hash].receivedAt == 0 {
		bm.download[hash].receivedAt = time.Now().UnixNano()
	}

	//on recupere le dernier block de la chain
	lastChainBlock := bm.chain.GetLastBlock()

	//SYSTEM ERROR - highly possible
	//si le nouveau block ne comporte pas le hash du dernier block de la chain dans son header
	if bytes.Compare(lastChainBlock.GetHash(), new.Header.HashPrevBlock) != 0 {
		//on recupere le temps moyen que met le noeud a téléchargé un block
		//valeur recuperé en nanosecond
		averageTimeToDownloadBlock := GetAverageTimeToDownloadABlock(bm.download)
		//si il n'y a pas encore assez de data permettant de determiner le temps moyen
		if averageTimeToDownloadBlock == 0 {
			//500 MS
			averageTimeToDownloadBlock = 1000000 * 500
		}
		if bm.download[hash].nbTry == 0 {
			bm.download[hash].unConfirmedBlock = new
		}
		//on ajoute la date du nouvel essai pour ajoute le block à la chaine.
		//nouvelle date = date actuel + temps moyen pour dl un block
		bm.download[hash].nbTry += 1
		bm.download[hash].timeToRetry = time.Now().Add(time.Nanosecond * time.Duration(averageTimeToDownloadBlock)).UnixNano()
		cleared := bm.ClearPendingBlockIfNecessary(hash)
		if cleared == false {
			go func() {
				//on sleep dans une goroutine en attendant le nouvel essai
				time.Sleep(time.Nanosecond * time.Duration(averageTimeToDownloadBlock))
				bm.BlockDownloaded(new, s)
			}()
		}
		return
	}

	//on check le nouveau block
	//si il y a une erreur, on supprime le block du manager, le block est invalide
	err := bm.chain.CheckNewBlock(new)
	if err != nil {
		bm.Log(false, "wrong new block informations")
		delete(bm.download, hash)
		return
	}

	//on ajoute le block à la chain
	err = bm.chain.AddBlock(new)
	if err == nil {
		//met a jour le nouveau tip du mining manager
		s.MiningManager.UpdateTip(new.GetHash())
		//Si le noeud est en cours de minage
		if s.MiningManager.IsMining() == true {
			s.newBlock <- new
		} else if s.mining == true && s.IsNodeAbleToMine() {
			go s.Mining()
		}
		if bm.chain.Height == s.newFetchAtHeight {
			s.askNewBlock(nil, -1)
		}

		bm.Log(true, fmt.Sprintf("block %d - %s successfully added on chain\n", bm.chain.Height, hash))
		bm.download[hash].block = new
		bm.download[hash].addedAt = time.Now().UnixNano()
	} else {
		bm.Log(true, fmt.Sprintf("block %s hasn't been added on chain next to this error: %s\n", hash, err.Error()))
		delete(bm.download, hash)
	}
}

//Cette fonction est appelé lorsque l'on commence à télécharger un block depuis un pair
func (bm *blockManager) StartDownloadBlock(hash string, sp *serverPeer, expectedHeight int64) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	bm.download[hash] = &DownloadInformations{
		sp:             sp,
		start:          time.Now().UnixNano(),
		expectedHeight: expectedHeight,
	}
	bm.Log(true, fmt.Sprintf("Start downloading %s from %s", hash, sp.GetAddr()))
}

func (bm *blockManager) Log(printTime bool, c ...interface{}) {
	if bm.log == false {
		return
	}
	fmt.Print("Block Manager: ")
	if bm.log == true {
		if printTime == true {
			fmt.Print(time.Now().Format("15:04:05.000000"))
			fmt.Print(" ")
		}
		for _, seq := range c {
			fmt.Print(seq)
			fmt.Print(" ")
		}
		fmt.Print("\n")
	}
}

func (bm *blockManager) GetDownloadedBlock(hash []byte) *twayutil.Block {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	hashString := hex.EncodeToString(hash)
	d, exist := bm.download[hashString]
	if exist {
		return d.block
	}
	return nil
}

func (bm *blockManager) IsContainOrphanBlock() bool {
	bm.mu.Lock()
	cpy := bm.download
	bm.mu.Unlock()
	var i = 0
	for _, d := range cpy {
		if d.nbTry > 0 && d.block == nil {
			i++
		}
		if i >= conf.MaxBlockPerMsg {
			return true
		}
	}
	return false
}

func (bm *blockManager) GetUncomfirmedBlocks() listDownloadInformation {
	var ret listDownloadInformation
	bm.mu.Lock()
	cpy := bm.download
	bm.mu.Unlock()
	for _, d := range cpy {
		if d.block == nil && d.unConfirmedBlock != nil {
			ret = append(ret, *d)
		}
	}
	return ret
}

func (bm *blockManager) ClearPendingBlockIfNecessary(hash string) bool {
	di, exist := bm.download[hash]
	if exist == false {
		return false
	}
	twoMinutesBack := time.Now().Add(time.Minute * -2)
	unconfirmReceivedAt := time.Unix(0, di.receivedAt)
	if unconfirmReceivedAt.Before(twoMinutesBack) {
		delete(bm.download, hash)
		return true
	}
	return false
}

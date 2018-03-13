package mining

import (
	b "tway/blockchain"
	"tway/twayutil"
	"tway/wallet"
	"tway/util"
	"math/big"
	"time"
	"fmt"
	"encoding/hex"
	"math"
	"bytes"
	"sync"
)

var mempool []twayutil.Transaction
const maxNonce = math.MaxInt64

type MiningManager struct {
	NewBlock 		chan *twayutil.Block
	is_mining 		bool
	start			time.Time
	tip				[]byte
	quit 			chan int
	log				bool
	chain 			*b.Blockchain	

	mu 				sync.Mutex
	HistoryMined 	map[string]*twayutil.Block	
}

func NewMiningManager(tip []byte, log bool, chain *b.Blockchain) *MiningManager {
	return &MiningManager{
		NewBlock: make(chan *twayutil.Block),
		HistoryMined: make(map[string]*twayutil.Block),
		tip: tip,
		quit: make(chan int),
		log: log,
		chain: chain,
	}
}

func (mm *MiningManager) IsMining() bool {
	return mm.is_mining
}

func (mm *MiningManager) UpdateTip(newTip []byte){
	mm.tip = newTip
}

func (mm *MiningManager) run(pow *b.Pow, newBlock chan *twayutil.Block, quit chan int){
	var hashInt big.Int
	var hash []byte
	nonce := 0
	var stopMining = false

	go func(){
		for nonce < maxNonce {
			data := pow.PrepareData(util.EncodeInt(nonce))
			hash = util.Sha256(data)
			hashInt.SetBytes(hash[:])
			if hashInt.Cmp(pow.Target) == -1 {
				pow.Block.Header.Nonce = util.EncodeInt(nonce)
				pow.Block.Size = util.EncodeInt(int(pow.Block.GetSize()))

				go func(){
					//signal au serveur qu'un nouveau block a été miné
					//le serveur traitera le block via le block manager
					//et informera le mining manager qu'un nouveau block a été ajouté à la chain
					mm.NewBlock <- pow.Block
					//signal que le minage de ce block est terminé.
					//pour passer au block suivant
					newBlock <- pow.Block
				}()
				return
			}
			if stopMining == true {
				return 
			}
			nonce++
		}
	}()
	for {
		select {
			case new := <-newBlock:
				if bytes.Compare(pow.Block.Header.HashPrevBlock, new.Header.HashPrevBlock) == 0 {
					stopMining = true
					mm.tip = new.GetHash()
					return
				}
			case <-mm.quit:
				quit <- 1
				stopMining = true
				return
		}
	}
}

func (mm *MiningManager) StartMining(newBlock chan *twayutil.Block, tip []byte){
	if mm.is_mining == true {
		return
	}	
	mm.is_mining = true
	quit := make(chan int)
	var stop = false
	go func(){
		<-quit
		stop = true
	}()
	mm.Log(false, "started")
	for stop == false {
		txs := mempool
		_, _, fees := b.GetTotalAmounts(txs)
		time.Sleep(100 * time.Millisecond)
		block := twayutil.NewBlock(txs, mm.tip, wallet.NewMiningWallet(), fees, mm.chain.GetNewBits())
		//Créer une target de proof of work
		pow := b.NewProofOfWork(block)
		mm.run(pow, newBlock, quit)	
	}
	mm.Log(false, "stopped")
	mm.is_mining = false
}

func (mm *MiningManager) Stop(){
	mm.quit <- 1
}

func (mm *MiningManager) Log(printTime bool, c... interface{}){
	if mm.log == false {
		return 
	}
	fmt.Print("Mining: ")
	if (mm.log == true){
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

func (mm *MiningManager) AddToHistoryMined(newBlock *twayutil.Block) {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	hash := hex.EncodeToString(newBlock.GetHash())
	mm.HistoryMined[hash] = newBlock
}

func (mm *MiningManager) GetBlockInHistory(hash []byte) *twayutil.Block{
	mm.mu.Lock()
	defer mm.mu.Unlock()

	hashString := hex.EncodeToString(hash)
	b, exist := mm.HistoryMined[hashString]
	if exist {
		return b
	}
	return nil
}
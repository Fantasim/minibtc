package mining

import (
	b "tway/blockchain"
	"tway/twayutil"
	"tway/wallet"
	"tway/util"
	"fmt"
	"math/big"
	"time"
	"math"
	"bytes"
)

var mempool []twayutil.Transaction
const maxNonce = math.MaxInt64

type MiningManager struct {
	NewBlock 	chan *twayutil.Block
	is_mining 	bool
	start		time.Time
	tip			[]byte
}

func NewMiningManager(tip []byte) *MiningManager {
	return &MiningManager{
		NewBlock: make(chan *twayutil.Block),
		tip: tip,
	}
}

func (mm *MiningManager) IsMining() bool {
	return mm.is_mining
}

func (mm *MiningManager) UpdateTip(newTip []byte){
	mm.tip = newTip
}

func (mm *MiningManager) run(pow *b.Pow, newBlock chan *twayutil.Block){
	var hashInt big.Int
	var hash []byte
	nonce := 0
	var stopMining = false

	go func(){
		mm.is_mining = true 
		for nonce < maxNonce {
			if nonce % 100000 == 0 {
			}
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
					i := 0;
					for {
						if bytes.Compare(mm.tip, pow.Block.Header.HashPrevBlock) == 0 {
							break;
						}
						time.Sleep(time.Microsecond * 10000)
						i++
					}
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
		new := <-newBlock
		if bytes.Compare(pow.Block.Header.HashPrevBlock, new.Header.HashPrevBlock) == 0 {
			stopMining = true
			mm.tip = new.GetHash()
			return
		}
	}
}

func (mm *MiningManager) StartMining(newBlock chan *twayutil.Block, tip []byte){
	mm.start = time.Now()
	fmt.Println("[MINING] START")
	for {
		txs := mempool
		_, _, fees := b.GetTotalAmounts(txs)
		block := twayutil.NewBlock(txs, mm.tip, wallet.RandomWallet().PublicKey, fees)
		//Créer une target de proof of work
		pow := b.NewProofOfWork(block)
		mm.run(pow, newBlock)
	}
	mm.is_mining = false
}
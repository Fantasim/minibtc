package mempool

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
	"tway/blockchain"
	"tway/twayutil"
)

var Mempool = NewMempool()

type listDownloadInformation []DownloadInformations

type TxPool struct {
	pool     sync.Map
	download sync.Map
	log      bool
}

type DownloadInformations struct {
	start      int64 //ns
	receivedAt int64 //ns
	addedAt    int64 //ns
}

func (tp *TxPool) StartDownloadTx(hash []byte) error {
	if _, ok := tp.download.Load(hash); ok == true {
		return errors.New("a download order has already been provided")
	}
	d := &DownloadInformations{
		start: time.Now().UnixNano(),
	}
	tp.download.Store(hash, d)
	return nil
}

func (tp *TxPool) RemoveDownloadInformation(hash []byte) {
	tp.download.Delete(hash)
}

func NewMempool() *TxPool {
	tp := &TxPool{
		pool:     sync.Map{},
		download: sync.Map{},
		log:      true,
	}
	return tp
}

func (tp *TxPool) CheckIfTxInputsAlreadyUsedInMempool(tx *twayutil.Transaction) error {
	TXs := tp.PoolToTxSlice()

	for _, poolTX := range TXs {
		for _, in := range poolTX.Inputs {
			for _, txInputs := range tx.Inputs {
				if bytes.Compare(txInputs.PrevTransactionHash, in.PrevTransactionHash) == 0 && bytes.Compare(txInputs.Vout, in.Vout) == 0 {
					return errors.New("An input of this tx is already spent in mempool")
				}
			}
		}
	}
	return nil
}

func (tp *TxPool) AddTx(tx *twayutil.Transaction) error {

	hash := hex.EncodeToString(tx.GetHash())
	/*diInterface, _ := tp.download.Load(tx.GetHash())
	if diInterface == nil {
		return errors.New("any download order has been provided for this tx.")
	}
	di := diInterface.(*DownloadInformations)
	di.receivedAt = time.Now().UnixNano()
	*/
	if res := tp.GetTx(hash); res != nil {
		return errors.New("this tx already exist in mempool")
	}
	if err := tp.CheckIfTxInputsAlreadyUsedInMempool(tx); err != nil {
		return err
	}
	if err := blockchain.CheckIfTxIsCorrect(tx); err != nil {
		return err
	}
	//go func() {
	tp.pool.Store(hex.EncodeToString(tx.GetHash()), tx)
	//di.addedAt = time.Now().UnixNano()
	//tp.download.Store(tx.GetHash(), di)
	//}()
	tp.Log(false, hex.EncodeToString(tx.GetHash()), " added")
	return nil
}

func (tp *TxPool) GetTx(hash string) *twayutil.Transaction {
	val, exist := tp.pool.Load(hash)
	if exist == false {
		return nil
	}
	ret := val.(*twayutil.Transaction)
	return ret
}

func (tp *TxPool) RemoveTxListIfExist(txs []twayutil.Transaction) {
	for _, tx := range txs {
		tp.RemoveTx(hex.EncodeToString(tx.GetHash()))
	}
}

func (tp *TxPool) RemoveTx(hash string) error {
	tp.pool.Delete(hash)
	return nil
}

func (tp *TxPool) PoolToTxSlice() []twayutil.Transaction {
	var ret []twayutil.Transaction
	tp.pool.Range(func(key, val interface{}) bool {
		tx := val.(*twayutil.Transaction)
		ret = append(ret, *tx)
		return true
	})
	return ret
}

func (tp *TxPool) Log(printTime bool, c ...interface{}) {
	if tp.log == false {
		return
	}
	fmt.Print("Mempool: ")
	if tp.log == true {
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

package blockchain

import (
	twayutil "tway/twayutil"
	conf "tway/config"
	util "tway/util"
	"bytes"
	"github.com/boltdb/bolt"
	"time"
	"strconv"
	"errors"
	"fmt"

)

func (b *Blockchain) GetBlockHeight(blockHash []byte) int {
	be := NewExplorer()
	var i = 0
	for {
		bl := be.Next();
		if bytes.Compare(bl.GetHash(), blockHash) == 0 {
			return BC.Height - i
		}
		if bl == nil {
			return -1
		}
		i++
	}
	return i
}

func (b *Blockchain) GetBlockByHash(hash []byte) (*twayutil.Block, int) {
	var block *twayutil.Block
	
	db := b.DB

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCK_BUCKET))
		encodedBlock := b.Get(hash)
		if len(encodedBlock) == 0 {
			return errors.New("Block doesn't exist")
		}
		block = twayutil.DeserializeBlock(encodedBlock)
		return nil
	})
	if err != nil {
		return nil, -1
	}
	return block, b.GetBlockHeight(block.GetHash())
}

func GenesisBlock(pubKey []byte) *twayutil.Block {
	//créer une nouvelle transaction coinbase
	tx := twayutil.NewCoinbaseTx(pubKey, 0)

	//créer un nouveau block et ajoute la structure
	block := &twayutil.Block{
		Transactions: []twayutil.Transaction{tx},
		Counter: 1,
	}

	//recupère le merkle root de la liste de transaction
	//contenant uniquement la transaction coinbase
	HashMerkleRoot := util.NewMerkleTree([][]byte{tx.Serialize()}).RootNode.Data

	//Créer le header du block
	header := twayutil.BlockHeader{
		Version: []byte{conf.VERSION},
		HashPrevBlock: conf.GENESIS_BLOCK_PREVHASH,
		HashMerkleRoot: HashMerkleRoot,
		Time:  util.EncodeInt(int(time.Now().Unix())),
		Bits:  util.EncodeInt(1),
	}
	block.Header = header
	MineBlock(block)
	return block
}

func MineBlock(b *twayutil.Block) error {
	//Créer une target de proof of work
	pow := NewProofOfWork(b)
	//cherche le nonce correspondant à la target
	nonce, _, err := pow.Run()
	if err != nil {
		return err
	}
	//ajoute le nonce au header
	b.Header.Nonce = util.EncodeInt(nonce)
	//ajoute la taille total du block
	b.Size = util.EncodeInt(int(b.GetSize()))
	return nil
}

func GetTotalFees(list []twayutil.Transaction) int {
	var total_fees = 0
	for _, tx := range list {
		total_inputs, total_outputs, fees := GetAmounts(tx.GetHash())
		if total_outputs > total_inputs {
			fmt.Println("Total outputs is greater than total_inputs. dectecting cheat try.")
			return -1
		}
		total_fees += fees
	}
	return total_fees
}

func (b *Blockchain) GetNBlocksNextToHeight(height int, max int) map[string]*twayutil.Block {
	var list = make(map[string]*twayutil.Block)
	
	be := NewExplorer()
	for i := height; i < b.Height; i++ {
		block := be.Next()
		if len(list) == conf.MaxBlockPerMsg || block == nil || len(list) == max {
			break
		}
		if (b.Height - i) < conf.MaxBlockPerMsg {
			list[strconv.Itoa(i)] = block
		}
	}
	return list
}
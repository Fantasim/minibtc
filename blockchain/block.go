package blockchain

import (
	"letsgo/wallet"
	"unsafe"
	"bytes"
	"github.com/boltdb/bolt"
	"encoding/gob"
	"time"
	"letsgo/util"
	"log"
	"fmt"
)

//HashPrevBlock du block genèse
var GENESIS_BLOCK_PREVHASH = []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}

//Structure d'un Block
type Block struct {
	Size []byte //taille du block en octet
	Header BlockHeader //Header du block
	Counter uint //nombre de transaction
	Transactions []Transaction //liste de transaction
}

//Structure du header d'un block
type BlockHeader struct{
	Version []byte //version du noeud créateur du block
	HashPrevBlock []byte //hash du dernier block de la blockchain
	HashMerkleRoot []byte //merkleroot des transactions du block
	Time []byte //time unix de la création du block 
	Bits []byte //niveau de difficulté de minage
	Nonce []byte //nombre d'iteration nécéssaire pour trouver la solution de minage
}

//Créer un block genese
func NewGenesisBlock(address string) Block {

	//récupère la clé public liée à l'address
	pubKey := wallet.GetPubKeyFromAddress(address)
	//créer une nouvelle transaction coinbase
	tx := NewCoinbaseTx(pubKey)

	//créer un nouveau block et ajoute la structure
	block := Block{
		Transactions: []Transaction{tx},
		Counter: 1,
	}

	//recupère le merkle root de la liste de transaction
	//contenant uniquement la transaction coinbase
	HashMerkleRoot := NewMerkleTree([][]byte{tx.Serialize()}).RootNode.Data

	//Créer le header du block
	header := BlockHeader{
		Version: []byte{VERSION},
		HashPrevBlock: GENESIS_BLOCK_PREVHASH,
		HashMerkleRoot: HashMerkleRoot,
		Time:  util.EncodeInt(int(time.Now().Unix())),
		Bits:  util.EncodeInt(1),
	}
	block.Header = header
	//Créer une target de proof of work
	pow := NewProofOfWork(&block)
	//cherche le nonce correspondant à la target
	nonce, _, err := pow.Run()
	if err != nil {
		log.Panic(err)
	}
	//ajoute le nonce au header
	block.Header.Nonce = util.EncodeInt(nonce)
	//ajoute la taille total du block
	block.Size = util.EncodeInt(int(unsafe.Sizeof(block)))
	return block
}

//Retourne le hash d'un block
func (b *Block) GetHash() []byte {
	hash := bytes.Join(
		[][]byte{
			b.Header.HashPrevBlock,
			b.Header.HashMerkleRoot,
			b.Header.Time,
			b.Header.Bits,
			b.Header.Nonce,
			b.Size,
		},
		[]byte{},
	)

	return util.Sha256(hash)
}

//Serialize un block 
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

//Deserialize un block
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

func NewBlock(txs []Transaction, prevBlockHash []byte) *Block{
	block := &Block{}
	
	//Récupère un wallet aléatoire vers qui envoyer la transaction coinbase
	w := wallet.RandomWallet()

	fmt.Println(string(w.GetAddress()))
	//Créer une transaction coinbase
	coinbaseTx := NewCoinbaseTx(w.PublicKey)

	//Prepend la transaction coinbase à liste de transaction
	txs = append([]Transaction{coinbaseTx}, txs...)

	block.Transactions = txs
	block.Counter = uint(len(txs))

	//[]Transaction to [][]byte
	txsDoubleByteArray := TransactionToByteDoubleArray(txs)

	//recupère le merkle root de la liste de transaction
	HashMerkleRoot := NewMerkleTree(txsDoubleByteArray).RootNode.Data

	//Header du block
	header := BlockHeader{
		Version: []byte{VERSION},
		HashPrevBlock: prevBlockHash,
		HashMerkleRoot: HashMerkleRoot,
		Time:  util.EncodeInt(int(time.Now().Unix())),
		Bits:  util.EncodeInt(1),
	}
	block.Header = header


	//Créer une target de proof of work
	pow := NewProofOfWork(block)
	//cherche le nonce correspondant à la target
	nonce, _, err := pow.Run()
	if err != nil {
		log.Panic(err)
	}
	//ajoute le nonce au header
	block.Header.Nonce = util.EncodeInt(nonce)
	//ajoute la taille total du block
	block.Size = util.EncodeInt(int(unsafe.Sizeof(block)))
	return block
}

func (block *Block) GetBlockHeight() int {
	be := NewExplorer()
	var i = 0
	for {
		bl := be.Next();
		if bytes.Compare(bl.GetHash(), block.GetHash()) == 0 {
			return BC_HEIGHT - i
		}
		if bl == nil {
			return -1
		}
		i++
	}
	return i
}

func GetBlockByHash(hash []byte) (*Block, error) {
	var block *Block
	
	db := BC.DB

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCK_BUCKET))
		encodedBlock := b.Get(hash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return block, nil
}
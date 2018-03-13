package twayutil

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"time"
	"tway/util"
	conf "tway/config"
	"log"
)

type Blocks []Blocks

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
func (bl *Block) Serialize() []byte {
	b, err := json.Marshal(bl)
	if err != nil {
        log.Panic(err)
	}
	bu := new(bytes.Buffer)
	enc := gob.NewEncoder(bu)
	err = enc.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return bu.Bytes()
}

func GetListBlocksHashFromSlice(list []*Block) [][]byte {
	var ret [][]byte
	for _, b := range list {
		ret = append(ret, b.GetHash())
	}
	return ret
}

func GetListBlocksHashFromMap(list map[string]*Block) [][]byte {
	var ret [][]byte
	for _, b := range list {
		ret = append(ret, b.GetHash())
	}
	return ret
}

func (b *Block) GetSize() uint64 {
	 return 0
}

//Deserialize un block
func DeserializeBlock(data []byte) *Block {
	bl := new(Block)
	var dataByte []byte

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&dataByte)
	if err != nil {
		log.Panic(err)
	}
	json.Unmarshal(dataByte, bl)
	return bl
}

func GetMerkleHash(txs []Transaction) []byte {
	//[]Transaction to [][]byte
	txsDoubleByteArray := TransactionToByteDoubleArray(txs)

	//recupère le merkle root de la liste de transaction
	return util.NewMerkleTree(txsDoubleByteArray).RootNode.Data
}

func NewBlock(txs []Transaction, prevBlockHash []byte, pubKeyCoinbase []byte, total_fees int, bits int64) *Block{
	block := &Block{}
	//Récupère un wallet aléatoire vers qui envoyer la transaction coinbase

	//Créer une transaction coinbase
	coinbaseTx := NewCoinbaseTx(pubKeyCoinbase, total_fees)

	//Prepend la transaction coinbase à liste de transaction
	txs = append([]Transaction{coinbaseTx}, txs...)

	block.Transactions = txs
	block.Counter = uint(len(txs))
	//Header du block
	header := BlockHeader{
		Version: []byte{conf.VERSION},
		HashPrevBlock: prevBlockHash,
		HashMerkleRoot: GetMerkleHash(txs),
		Time:  util.EncodeInt(int(time.Now().Unix())),
		Bits:  util.EncodeInt(int(bits)),
	}
	block.Header = header

	return block
}


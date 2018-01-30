package wire

import (
	"bytes"
	"encoding/gob"
	"time"
	"tway/util"
	conf "tway/config"
	"log"
)

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
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

func (b *Block) GetSize() uint64 {
	 return 0
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

func NewBlock(txs []Transaction, prevBlockHash []byte, pubKeyCoinbase []byte, total_fees int) *Block{
	block := &Block{}
	//Récupère un wallet aléatoire vers qui envoyer la transaction coinbase

	//Créer une transaction coinbase
	coinbaseTx := NewCoinbaseTx(pubKeyCoinbase, total_fees)

	//Prepend la transaction coinbase à liste de transaction
	txs = append([]Transaction{coinbaseTx}, txs...)

	block.Transactions = txs
	block.Counter = uint(len(txs))

	//[]Transaction to [][]byte
	txsDoubleByteArray := TransactionToByteDoubleArray(txs)

	//recupère le merkle root de la liste de transaction
	HashMerkleRoot := util.NewMerkleTree(txsDoubleByteArray).RootNode.Data

	//Header du block
	header := BlockHeader{
		Version: []byte{conf.VERSION},
		HashPrevBlock: prevBlockHash,
		HashMerkleRoot: HashMerkleRoot,
		Time:  util.EncodeInt(int(time.Now().Unix())),
		Bits:  util.EncodeInt(1),
	}
	block.Header = header

	return block
}


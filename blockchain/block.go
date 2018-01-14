package blockchain

import (
	//"strconv"
	//"letsgo/wallet"
	//"unsafe"
	"bytes"
	"encoding/gob"
	//"time"
	"letsgo/util"
	"log"
	"encoding/hex"
)

var GENESIS_BLOCK_PREVHASH = []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}

type Block struct {
	Size []byte
	Header BlockHeader
	Counter uint
	Transactions []Transaction
}

type BlockHeader struct{
	Version []byte
	HashPrevBlock []byte
	HashMerkleRoot []byte
	Time []byte
	Bits []byte
	Nonce []byte
}

/*
func NewGenesisBlock() Block {

	pubKey := wallet.GetPubKeyFromAddress("16caHAfC5FpWWtmXTqphtQyRUXN2DgorJ3")
	signature, _  := wallet.SignPrivateKey(pubKey)
	tx := NewCoinbaseTx(pubKey, signature)

	block := Block{
		Transactions: []Transaction{tx},
		Counter: 1,
	}

	HashMerkleRoot := NewMerkleTree([][]byte{tx.Serialize()}).RootNode.Data

	header := BlockHeader{
		Version: []byte{VERSION},
		HashPrevBlock: GENESIS_BLOCK_PREVHASH,
		HashMerkleRoot: HashMerkleRoot,
		Time:  []byte(strconv.Itoa(int(time.Now().Unix()))),
		Bits: []byte("1"),
	}
	block.Header = header
	pow := NewProofOfWork(&block)
	nonce, _, err := pow.Run()
	if err != nil {
		log.Panic(err)
	}
	block.Header.Nonce = util.IntToArrayByte(nonce)
	block.Size = util.IntToArrayByte(int(unsafe.Sizeof(block)))

	return block
} */


func NewGenesisBlock() Block {
	blockStr := "46ff8103010105426c6f636b01ff82000104010453697a65010a00010648656164657201ff84000107436f756e746572010600010c5472616e73616374696f6e7301ff9200000066ff830301010b426c6f636b48656164657201ff84000106010756657273696f6e010a00010d4861736850726576426c6f636b010a00010e486173684d65726b6c65526f6f74010a00010454696d65010a00010442697473010a0001054e6f6e6365010a00000027ff91020101185b5d626c6f636b636861696e2e5472616e73616374696f6e01ff920001ff86000068ff850301010b5472616e73616374696f6e01ff86000106010756657273696f6e010a000109496e436f756e746572010a000106496e7075747301ff8c00010a4f7574436f756e746572010a0001074f75747075747301ff900001084c6f636b54696d65010a00000021ff8b020101125b5d626c6f636b636861696e2e496e70757401ff8c0001ff88000055ff8703010105496e70757401ff880001040113507265765472616e73616374696f6e48617368010a000104566f7574010a00010d5478496e5363726970744c656e010a00010953637269707453696701ff8a00000017ff89020101095b5d5b5d75696e743801ff8a00010a000022ff8f020101135b5d626c6f636b636861696e2e4f757470757401ff900001ff8e000043ff8d030101064f757470757401ff8e000103010556616c7565010a00010e54785363726970744c656e677468010a00010c5363726970745075624b657901ff8a000000fe011eff8201033230300101010001200000000000000000000000000000000000000000000000000000000000000000012097fd53e6f4677253826c6f9e772a75ac7049962bb73b7dba67f403c7e2bad625010a313531353935353134360101310107323635323432340001010101010100010131010102022d3101023635010240b39cc446ceac983bd1d9a167920affb9f264fc8e44e4266efcda46fe32a13d73cbcf337561b5434e71ba64330b3e1b9b880118f6944e6114c87cb2ac8bd6513b01ac000101310101010835303030303030300102363401014028f06c3debd5047a16336473a58f413a2577fcbbaeb8a31799fa0c50400f1cfbb41284e7697bfdb630acabdfbbec9e3db224f14e18acfcd6e9689ca7f66bd62e000101000000"	
	b, err := hex.DecodeString(blockStr)
	if err != nil {
		log.Panic(err)
	}
	return *DeserializeBlock(b)
}

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
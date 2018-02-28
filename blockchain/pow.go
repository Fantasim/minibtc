package blockchain

import ( 
	"math/big"
	"math"
	"bytes"
	"tway/util"
	"tway/twayutil"
	"time"
	"fmt"
)

var maxNonce = math.MaxInt64
const targetBits = 23

type Pow struct {
	Block  *twayutil.Block
	Target *big.Int
}

// NewProofOfWork builds and returns a ProofOfWork
func NewProofOfWork(b *twayutil.Block) *Pow {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &Pow{b, target}

	return pow
}

func (pow *Pow) PrepareData(nonce []byte) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.Header.HashPrevBlock,
			pow.Block.Header.HashMerkleRoot,
			pow.Block.Header.Time,
			pow.Block.Header.Bits,
			nonce,
		},
		[]byte{},
	)

	return data
}

//Cherche un hash inférieur à la target (mine)
func (pow *Pow) Run() (int, []byte, error) {
	var hashInt big.Int
	var hash []byte

	nonce := 0
	start := time.Now()

	fmt.Println("mining...")
	for nonce < maxNonce {
		data := pow.PrepareData(util.EncodeInt(nonce))
		hash = util.Sha256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.Target) == -1 {
			break
		}
		nonce++
	}
	fmt.Println("nonce =", nonce)
	fmt.Println("Solution found after: ", time.Now().Sub(start), "\n")
	return nonce, hash[:], nil
}

//Verifie que la prof of work est bien validé par la règle imposée.
func (pow *Pow) Validate() bool {
	var hashInt big.Int

	data := pow.PrepareData(pow.Block.Header.Nonce)
	hash := util.Sha256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.Target) == -1
}
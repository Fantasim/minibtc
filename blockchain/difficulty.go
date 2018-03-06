package blockchain

import (
	"math/big"
	conf "tway/config"
	"tway/util"
	"fmt"
)

func MaxInt() *big.Int {
	max := big.NewInt(1)
	max.Lsh(max, 256)
	return max
}

func GetInitialTarget() *big.Int {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	return target
}

func CalcNextTarget(n *big.Int) *big.Int {
    i := GetInitialTarget()
	t := i.Div(i, n)
    return t
}

func (b *Blockchain) GetNewBits() int64 {
	//si la hauteur de chain + 1 est inférieur à l'interval de block, 
	//tous lesquels la difficulté est recalculé 
	if b.Height + 1 < conf.NEW_DIFFICULTY_EACH_N_BLOCK {
		return int64(util.DecodeInt(b.GetGenesisBlock().Header.Bits))
	}
	//dernier block de la chaine
	lastBlock := b.GetLastBlock()
	prevNBlock := b.GetBlockByHeight((b.Height + 1 + 1) - conf.NEW_DIFFICULTY_EACH_N_BLOCK)

	if (b.Height + 1) % conf.NEW_DIFFICULTY_EACH_N_BLOCK == 0 {
		genesisBlock := b.GetBlockByHeight(1)
	
		lastBlockTime := util.DecodeInt(lastBlock.Header.Time)
		genesisBlockTime := util.DecodeInt(genesisBlock.Header.Time)
	
		timeSinceBeginning := lastBlockTime - genesisBlockTime
		targetTime := conf.TARGET_TIME_BETWEEN_TWO_BLOCKS * b.Height

		diviseTargetBy := float64(targetTime) / float64(timeSinceBeginning)
		lastBits := util.DecodeInt(prevNBlock.Header.Bits)
		if lastBits == 1 && diviseTargetBy < 1 {
			return 1
		}

		newBits := float64(lastBits) * diviseTargetBy

		if diviseTargetBy > 1 {
			if newBits - float64(int(newBits)) > 0 {
				newBits += 1
			}
		}
		
		fmt.Println("LAST bits", lastBits)
		fmt.Println("NEW bits", int64(newBits))
		return int64(newBits)
	}
	return int64(util.DecodeInt(lastBlock.Header.Bits))
}

package util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
)

type nodeList []MerkleNode

// MerkleNode represent a Merkle tree node
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

func (nl nodeList) Print() {
	for _, n := range nl {
		fmt.Print(hex.EncodeToString(n.Data))
		fmt.Print(" ")
	}
	fmt.Print("\n")
}

// NewMerkleNode creates a new Merkle tree node
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	mNode := MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		mNode.Data = hash[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		mNode.Data = hash[:]
	}

	mNode.Left = left
	mNode.Right = right
	return &mNode
}

func GetMerkleRoot(data [][]byte) *MerkleNode {
	inc := 0
	if len(data)%2 != 0 {
		inc = 1
	}
	numberOfTwoPower := GetNumberOfTwoPower(len(data)+inc) - 1

	var nodes nodeList
	nodes = make(nodeList, int(math.Pow(2, float64(numberOfTwoPower))))

	var i int
	var d []byte
	for i, d = range data {
		nodes[i] = *NewMerkleNode(nil, nil, d)
	}

	if (i+1)%2 != 0 {
		nodes[i+1] = nodes[i]
	}

	for i := 0; i < numberOfTwoPower; i++ {
		tmpNodes := make(nodeList, int(math.Pow(2, float64(numberOfTwoPower-i-1))))
		var j = 0
		for j < len(nodes)/2 {
			if len(nodes[j*2].Data) != 0 {
				tmpNodes[j] = *NewMerkleNode(&nodes[j*2], &nodes[(j*2)+1], []byte{})
			} else {
				break
			}
			j++
		}
		if j%2 != 0 && len(tmpNodes) > 2 {
			tmpNodes[j] = tmpNodes[j-1]
		}
		nodes = tmpNodes
	}

	return &nodes[0]
}

func GetNumberOfTwoPower(n int) int {
	cpt := 0
	i := 1
	for n > i {
		cpt++
		i *= 2
	}
	if i >= n {
		cpt++
	}
	return cpt
}

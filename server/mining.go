package server

import (
	"bytes"
	"tway/twayutil"
	"encoding/hex"
	"fmt"
)

func (s *Server) HandleNewBlockMined(){
	for {
		new := <- s.MiningManager.NewBlock
		s.Log(true, "[MINING] block", hex.EncodeToString(new.GetHash()), "mined !")
		//on recupere le dernier block de la chain
		lastChainBlock := s.chain.GetLastBlock()
		err := s.chain.CheckNewBlock(new);
		if bytes.Compare(lastChainBlock.GetHash(), new.Header.HashPrevBlock) == 0 && err == nil {
			err = s.chain.AddBlock(new)
			if err == nil {
				s.Log(false, "[MINING] successfully added ON CHAIN")
				s.Log(false, "[MINING] Height:", s.chain.Height)
				s.Log(false, "[MINING] ", hex.EncodeToString(new.GetHash()))
				listBlockTmp := make([]*twayutil.Block, 1)
				listBlockTmp[0] = new
				list := twayutil.GetListBlocksHashFromSlice(listBlockTmp)
				percentageOfSuccess := s.BootstrapInv("block", list)
				if percentageOfSuccess == 0 {
					s.Log(false, "/!/ FAIL TO SEND BLOCK MINED")
					return
				}
			} else {
				s.Log(false, "/!/ FAIL TO ADD ON CHAIN BLOCK MINED")
			}
		} else {
			s.Log(false, "/!/ Block is not next to current TIP")
		}
	}
}

func (s *Server) Mining(){
	fmt.Println("server/mining.go -> mining start")
	s.MiningManager.StartMining(s.newBlock, s.chain.Tip)
	fmt.Println("server/mining.go -> mining stop")
}

func (s *Server) IsNodeAbleToMine() bool {
	list := s.GetListOfTrustedMainNode()
	var i = 0
	for _, p := range list {
		if p.GetLastBlock() <= int64(s.chain.Height) {
			i++
		}
	}
	return i == len(list)
}
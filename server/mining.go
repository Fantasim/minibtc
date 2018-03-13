package server

import (
	"bytes"
	"encoding/hex"
	"sync"
)

func (s *Server) HandleNewBlockMined(){
	mu := sync.Mutex{}
	for {
		mu.Lock()
		new := <- s.MiningManager.NewBlock

		copyNew := new
		hashNew := copyNew.GetHash()
		s.MiningManager.Log(true, "block", hex.EncodeToString(hashNew), "mined !")
		//on recupere le dernier block de la chain
		lastChainBlock := s.chain.GetLastBlock()
		if bytes.Compare(lastChainBlock.GetHash(), copyNew.Header.HashPrevBlock) == 0 {
			err := s.chain.AddBlock(copyNew)
			if err == nil {
				s.MiningManager.AddToHistoryMined(copyNew)
				s.MiningManager.Log(false, "successfully added ON CHAIN")

				var list [][]byte
	
				list = append(list, hashNew)
				percentageOfSuccess := s.BootstrapInv("block", list)
				if percentageOfSuccess == 0 {
					s.MiningManager.Log(false, "/!/ FAIL TO SEND BLOCK MINED")
				}
			} else {
				s.MiningManager.Log(false, "/!/ FAIL TO ADD ON CHAIN BLOCK MINED")
				s.MiningManager.Stop()
			}
		} else {
			s.MiningManager.Log(false, "/!/ Block is not next to current TIP")
		}
		mu.Unlock()
	}
}

func (s *Server) Mining(){
	s.MiningManager.StartMining(s.newBlock, s.chain.Tip)
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
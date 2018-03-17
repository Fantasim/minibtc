package server

import (
	"bytes"
	"encoding/hex"
	"sync"
)

//function executé dans une goroutine.
//Elle attend qu'un nouveau block a été miné.
//Dès lors qu'un nouveau block est miné, elle l'ajoute a la chain
//et l'envoie a tout le réseau.
func (s *Server) HandleNewBlockMined(){
	mu := sync.Mutex{}

	for {
		mu.Lock()
		//attend un nouveau block miné
		new := <- s.MiningManager.NewBlock

		copyNew := new
		hashNew := copyNew.GetHash()

		s.MiningManager.Log(true, "block", hex.EncodeToString(hashNew), "mined !")

		//on recupere le dernier block de la chain
		lastChainBlock := s.chain.GetLastBlock()
		//si le hashprev du nouveau block miné correspond au hash du dernier block de la chain
		if bytes.Compare(lastChainBlock.GetHash(), copyNew.Header.HashPrevBlock) == 0 {
			//on ajoute le block
			err := s.chain.AddBlock(copyNew)
			if err == nil {
				//on ajoute le block a l'historique des blocks minés
				s.MiningManager.AddToHistoryMined(copyNew)
				s.MiningManager.Log(false, "successfully added ON CHAIN")

				var list [][]byte
				list = append(list, hashNew)
				//on envoie le block au réseau
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

//Commence le minage 
func (s *Server) Mining(){
	s.MiningManager.StartMining(s.newBlock, s.chain.Tip)
}

//Retourne true si le noeud courant possède 
//la dernière hauteur de block
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
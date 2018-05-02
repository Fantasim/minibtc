package server

import (
	"log"
	"time"
	conf "tway/config"
	"tway/serverutil"
)

func (s *Server) NewVersion(addrTo *serverutil.NetAddress) *serverutil.MsgVersion {
	return &serverutil.MsgVersion{
		ProtocolVersion: conf.NodeVersion,
		Timestamp:       time.Now(),
		AddrReceiver:    addrTo,
		AddrSender:      s.ipStatus,
		LastBlock:       s.chain.Height,
	}
}

//Envoie une structure de la version de notre blockchain au noeud principal
func (s *Server) sendVersion(addrTo *serverutil.NetAddress) ([]byte, error) {

	addr := addrTo.String()
	s.Log(true, "Version sent to:", addrTo.String())

	version := s.NewVersion(addrTo)

	payload := gobEncode(*version)
	request := append(commandToBytes("version"), payload...)

	err := s.sendData(addrTo.String(), request)

	if err == nil {
		go s.HistoryManager.NewVersionHistory(version, true)
		p, _ := s.GetPeer(addr)
		p.VersionSent()
		s.AddPeer(p)
	}

	return request, err
}

//Recupère la version d'un noeud
func (s *Server) handleVersion(request []byte) {
	var payload serverutil.MsgVersion
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}

	go s.HistoryManager.NewVersionHistory(&payload, false)

	s.Log(false, "\n")
	s.Log(true, "Version received from :", payload.AddrSender.String())
	s.Log(false, "\t - Block height:", payload.LastBlock)
	s.Log(false, "\t - Version:", payload.ProtocolVersion, "\n")

	//établie les informations concernant le pair
	//envoie un verack et sa version si non fait.
	go func() {
		addr := payload.AddrSender.String()
		p, _ := s.GetPeer(addr)
		//set la hauteur de chain du pair
		p.SetLastBlock(int64(payload.LastBlock))
		//set la hauteur de chain du pair lorsque la premiere connexion avec le pair a été etablie
		p.SetStartingHeight(int64(payload.LastBlock))
		//incremente le nombre de bytes recu depuis ce pair
		p.IncreaseBytesReceived(uint64(len(request)))
		//indique que le pair a envoyé sa version au moins une fois
		p.HasSentVersion()
		s.AddPeer(p)
		if _, err := s.sendVerack(payload.AddrSender); err == nil {
			p, _ = s.GetPeer(addr)
			if p.IsVersionSent() == false {
				s.sendVersion(payload.AddrSender)
			}
		}
	}()

	go func() {
		//si la hauteur de chain du noeud courant est
		//inférieur a la hauteur de chain du block emetteur
		//de la version
		if s.chain.Height < payload.LastBlock {
			//arrêter le minage si il est en cours d'execution.
			if s.MiningManager.IsMining() == true {
				s.MiningManager.Stop()
			}
			p, _ := s.GetPeer(payload.AddrSender.String())
			//lui demander des blocks
			s.askNewBlock(p, payload.LastBlock)
		}
	}()

	if s.mining == true && payload.AddrSender.IsEqual(GetMainNode()) && s.chain.Height >= payload.LastBlock && s.MiningManager.IsMining() == false {
		go s.Mining()
	}
}

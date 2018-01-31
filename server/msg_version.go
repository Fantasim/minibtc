package server

import (
	"time"
	b "tway/blockchain"
	conf "tway/config"
	"log"
	"fmt"
)

type MsgVersion struct {
	// Version of the protocol the node is using.
	ProtocolVersion int32
	// Time the message was generated.  This is encoded as an int64 on the wire.
	Timestamp time.Time
	// Address of the remote peer.
	AddrReceiver *NetAddress
	// Address of the local peer.
	AddrSender *NetAddress
	// Last block seen by the generator of the version message.
	LastBlock int
}

//Envoie une structure de la version de notre blockchain au noeud principal
func sendVersion(addrTo *NetAddress) error {
	//recupere la hauteur de la blockchain
	bestHeight := b.BC.Height

	payload := gobEncode(MsgVersion{conf.NodeVersion, time.Now(), addrTo, Node, bestHeight})

	request := append(commandToBytes("version"), payload...)
	//on envoie au noeud principale une requete
	return sendData(addrTo, request)
}

//Recupère la version d'un noeud
func handleVersion(request []byte) {
	var payload MsgVersion

	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}
	fmt.Println("Version received from :", payload.AddrSender.String())
	fmt.Println("Block height:", payload.LastBlock, "\n")
	//recupere la hauteur du noeud courant
	bestHeight := b.BC.Height
	//recupere la hauteur du noeud envoyant une version
	foreignerBestHeight := payload.LastBlock

	//si le height courant est inférieur au height du noeud recepteur
	if bestHeight < foreignerBestHeight {
		//on lui envoie une demande des blocks qu'il a
		//sendGetBlocks(payload.AddrFrom)
	} else if bestHeight > foreignerBestHeight {
		//on lui envoie notre version
		sendVersion(payload.AddrSender)
	}

	//si le noeud envoyant n'est pas un noeud connu
	if !nodeIsKnown(payload.AddrSender) {
		//on ajoute le noeud d'envoie dans la liste des noeuds connus
		KnownNodes = append(KnownNodes, payload.AddrSender)
	}
}
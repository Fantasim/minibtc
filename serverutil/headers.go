package serverutil

import "tway/twayutil"

type MsgAskHeaders struct {
	// Address of the local peer.
	AddrSender *NetAddress
	// Address of the local peer.
	AddrReceiver *NetAddress

	Version      int32
	HeadHash     []byte //Hash de fin dans l'intervalle demandé
	StoppingHash []byte //Hash de début dans l'intervalle demandé
	Count        uint16 //longueur de l'intervalle demandé
}

type Header struct {
	Height int
	Hash   []byte
	Header twayutil.BlockHeader
}

type MsgHeaders struct {
	// Address of the local peer.
	AddrSender *NetAddress
	// Address of the local peer.
	AddrReceiver *NetAddress

	Version int32
	List    []Header
}

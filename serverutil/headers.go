package serverutil

import (
	"encoding/hex"
	"fmt"
	"tway/twayutil"
)

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

	Version          int32
	List             []Header
	GetHeadersOrigin *MsgAskHeaders
}

func (mah *MsgAskHeaders) String() string {
	return fmt.Sprintf("{AddrSender: %s, AddrReceiver: %s, Version: %d, HeadHash: %s, StoppingHash: %s, Count: %d}",
		mah.AddrSender.String(),
		mah.AddrReceiver.String(),
		mah.Version,
		hex.EncodeToString(mah.HeadHash),
		hex.EncodeToString(mah.StoppingHash),
		mah.Count,
	)
}

func (mh *MsgHeaders) String() string {
	a := fmt.Sprintf("{AddrSender: %s, AddrReceiver: %s, Version: %d, List:\n")
	for index, h := range mh.List {
		a += fmt.Sprintf("[%d] %s\n", index, h.String())
	}
	return a
}

func (h *Header) String() string {
	return fmt.Sprintf("{Height: %d, Hash: %s, Header: %s}", h.Height, hex.EncodeToString(h.Hash), h.Header.String())
}

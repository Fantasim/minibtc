package serverutil

import "fmt"

type MsgAskBlocks struct {
	// Address of the local peer.
	Addr  *NetAddress
	Range [2]int
}

type MsgBlock struct {
	// Address of the local peer.
	AddrSender *NetAddress
	// Address of the local peer.
	AddrReceiver *NetAddress

	Data []byte
}

func (mab *MsgAskBlocks) String() string {
	return fmt.Sprintf("{Addr: %s, Range: [%d:%d]}", mab.Addr.String(), mab.Range[0], mab.Range[1])
}

package serverutil

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

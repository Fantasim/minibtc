package serverutil

type MsgInv struct {
	// Address of the local peer.
	AddrSender *NetAddress
	// Address of the local peer.
	AddrReceiver *NetAddress
	Kind         string // "tx" || "block"
	List         [][]byte
}

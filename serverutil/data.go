package serverutil

type MsgGetData struct {
	// Address of the local peer.
	AddrSender *NetAddress
	// Address of the local peer.
	AddrReceiver *NetAddress

	ID   []byte //hash du block ou de la tx
	Kind string //"block" ou "tx"
}

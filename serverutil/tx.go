package serverutil

import (
	"tway/twayutil"
)

type MsgTx struct {
	// Address of the local peer.
	AddrSender *NetAddress
	// Address of the local peer.
	AddrReceiver *NetAddress
	Tx           *twayutil.Transaction
}

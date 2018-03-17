package serverutil

import "time"

type MsgVersion struct {
	// Version of the protocol the node is using.
	ProtocolVersion int32
	// Time the message was generated.  This is encoded as an int64 on the twayutil.
	Timestamp time.Time
	// Address of the remote peer.
	AddrReceiver *NetAddress
	// Address of the local peer.
	AddrSender *NetAddress
	// Last block seen by the generator of the version message.
	LastBlock int
}

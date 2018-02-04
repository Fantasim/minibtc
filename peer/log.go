package peer

import (
	"fmt"
)

func (p *Peer) Print(){
	fmt.Println()
	fmt.Println(" | ", p.addr, " | ")
	fmt.Println("- VerAck received:", p.verAckReceived)
	fmt.Println("- Version sent:", p.versionSent)
	fmt.Println("- Last block:", p.lastBlock)
	fmt.Println("- Bytes received:", p.GetBytesReceived())
	fmt.Println("- Bytes sent:", p.GetBytesSent())
	fmt.Println()
}
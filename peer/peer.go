package peer

import (
	"sync/atomic"
	"sync"
)

type Peer struct {
	// The following variables must only be used atomically.
	bytesReceived uint64
	bytesSent     uint64

	addr    string

	statsMtx sync.Mutex
	lastBlock int64
	startingHeight int64

	versionSent          bool
	verAckReceived       bool
}

func NewPeer(addr string) *Peer{
	p := &Peer{}
	p.addr = addr
	return p
}

func (p *Peer) IsConfirmed() bool {
	return p.verAckReceived && p.versionSent
}

func (p *Peer) IncreaseBytesReceived(n uint64) {
	atomic.AddUint64(&p.bytesReceived, n)
}

func (p *Peer) IncreaseBytesSent(n uint64) {
	atomic.AddUint64(&p.bytesSent, n)
}

func (p *Peer) SetLastBlock(last int64){
	p.statsMtx.Lock()
	defer p.statsMtx.Unlock()
	p.lastBlock = last
}

func (p *Peer) SetStartingHeight(start int64){
	p.statsMtx.Lock()
	defer p.statsMtx.Unlock()
	if p.startingHeight == 0 {
		p.startingHeight = start	
	}
}

func (p *Peer) GetLastBlock() int64 {
	return p.lastBlock
}

func (p *Peer) GetAddr() string {
	return p.addr
}

func (p *Peer) VerAckReceived() {
	p.verAckReceived = true
}

func (p *Peer) VersionSent() {
	p.versionSent = true
}
package peer

import (
	"sync"
	"sync/atomic"
	"time"
)

type Peer struct {
	// The following variables must only be used atomically.
	bytesReceived uint64
	bytesSent     uint64

	addr string

	statsMtx       sync.Mutex
	lastBlock      int64
	startingHeight int64

	versionSent    bool
	verAckReceived bool
	hasSentVersion bool

	lastGetAddrTime      int64
	lastAskAddrTime      int64
	lastPingSentTime     int64
	lastPongReceivedTime int64
}

func NewPeer(addr string) *Peer {
	p := &Peer{}
	p.addr = addr
	return p
}

func (p *Peer) IncreaseBytesReceived(n uint64) {
	atomic.AddUint64(&p.bytesReceived, n)
}

func (p *Peer) IncreaseBytesSent(n uint64) {
	atomic.AddUint64(&p.bytesSent, n)
}

func (p *Peer) GetBytesReceived() uint64 {
	return p.bytesReceived
}

func (p *Peer) GetBytesSent() uint64 {
	return p.bytesSent
}

func (p *Peer) SetLastBlock(last int64) {
	p.statsMtx.Lock()
	defer p.statsMtx.Unlock()
	p.lastBlock = last
}

func (p *Peer) GotAddr() {
	p.statsMtx.Lock()
	defer p.statsMtx.Unlock()
	p.lastGetAddrTime = time.Now().UnixNano()
}

func (p *Peer) AskAddr() {
	p.statsMtx.Lock()
	defer p.statsMtx.Unlock()
	p.lastAskAddrTime = time.Now().UnixNano()
}

func (p *Peer) SetStartingHeight(start int64) {
	p.statsMtx.Lock()
	defer p.statsMtx.Unlock()
	if p.startingHeight == 0 {
		p.startingHeight = start
	}
}

func (p *Peer) GetLastAddrGetTime() int64 {
	p.statsMtx.Lock()
	defer p.statsMtx.Unlock()
	return p.lastGetAddrTime
}

func (p *Peer) GetLastBlock() int64 {
	return p.lastBlock
}

func (p *Peer) GetAddr() string {
	return p.addr
}

func (p *Peer) VerAckReceived() {
	p.statsMtx.Lock()
	defer p.statsMtx.Unlock()
	p.verAckReceived = true
}

func (p *Peer) VersionSent() {
	p.statsMtx.Lock()
	defer p.statsMtx.Unlock()
	p.versionSent = true
}

func (p *Peer) PingSent() {
	p.statsMtx.Lock()
	defer p.statsMtx.Unlock()
	p.lastPingSentTime = time.Now().UnixNano()
}

func (p *Peer) PongReceived() {
	p.statsMtx.Lock()
	defer p.statsMtx.Unlock()
	p.lastPongReceivedTime = time.Now().UnixNano()
}

func (p *Peer) HasSentVersion() {
	p.statsMtx.Lock()
	defer p.statsMtx.Unlock()
	p.hasSentVersion = true
}

func (p *Peer) IsVersionSent() bool {
	return p.versionSent
}

func (p *Peer) IsVerAckReceived() bool {
	return p.verAckReceived
}

func (p *Peer) HasHeSentVersion() bool {
	return p.hasSentVersion
}

func (p *Peer) GetLastPingSentTime() int64 {
	return p.lastPingSentTime
}

func (p *Peer) GetLastPongReceivedTime() int64 {
	return p.lastPongReceivedTime
}

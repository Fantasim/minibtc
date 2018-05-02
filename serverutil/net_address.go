package serverutil

import (
	"errors"
	"net"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidNetAddr = errors.New("provided net.Addr is not a net.TCPAddr")

// NetAddress defines information about a peer on the network including the time
// it was last seen, the services it supports, its IP address, and port.
type NetAddress struct {
	// Last time the address was seen.  This is, unfortunately, encoded as a
	// uint32 on the twayutil and therefore is limited to 2106.  This field is
	// not present in the decred version message (MsgVersion) nor was it
	// added until protocol version >= NetAddressTimeVersion.
	Timestamp time.Time

	// IP address of the peer.
	IP net.IP

	// Port the peer is using.  This is encoded in big endian on the twayutil
	// which differs from most everything else.
	Port uint16
}

func (addr *NetAddress) IsEqual(na *NetAddress) bool {
	if addr.IP.Equal(na.IP) && addr.Port == na.Port {
		return true
	}
	return false
}

func (addr *NetAddress) String() string {
	return addr.IP.String() + ":" + strconv.Itoa(int(addr.Port))
}

func NewNetAddress(addr net.Addr) (*NetAddress, error) {
	tcpAddr, ok := addr.(*net.TCPAddr)
	if !ok {
		return nil, ErrInvalidNetAddr
	}

	na := NewNetAddressIPPort(tcpAddr.IP, uint16(tcpAddr.Port))
	return na, nil
}

func NewNetAddressByString(addr string) (*NetAddress, error) {
	tcpAddr := strings.Split(addr, ":")
	if len(tcpAddr) == 2 {
		port, _ := strconv.Atoi(tcpAddr[1])
		na := NewNetAddressIPPort(net.ParseIP(tcpAddr[0]), uint16(port))
		return na, nil
	}
	return nil, ErrInvalidNetAddr
}

func NewNetAddressIPPort(ip net.IP, port uint16) *NetAddress {
	return NewNetAddressTimestamp(time.Now(), ip, port)
}

// NewNetAddressTimestamp returns a new NetAddress using the provided
// timestamp, IP, port, and supported services. The timestamp is rounded to
// single second precision.
func NewNetAddressTimestamp(
	timestamp time.Time, ip net.IP, port uint16) *NetAddress {
	// Limit the timestamp to one second precision since the protocol
	// doesn't support better.
	na := NetAddress{
		Timestamp: time.Unix(timestamp.Unix(), 0),
		IP:        ip,
		Port:      port,
	}
	return &na
}

package util

import (
	"errors"
	"net"
	"strings"
	"strconv"
)

func IpStringToBytes(ip string) []byte {
	var ret []byte
	sp := strings.Split(ip, ".")
	for _, k := range sp {
		n, _ := strconv.Atoi(k)
		ret = append(ret, byte(n))
	}
	return ret
}

func StringToNetIpAndPort(addr string) (net.IP, uint16) {
	t := strings.Split(addr, ":")
	if len(t) == 2{
		port, _ := strconv.Atoi(t[1])
		return IpStringToBytes(t[0]), uint16(port)
	}
	return IpStringToBytes(t[0]), 0
}

func GetIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip, nil
		}
	}
	return nil, errors.New("are you connected to the network?")
}

package gowip

import (
	"encoding/binary"
	"net"
)

// IpToUint transforms an ipv4 to a uint representation
// ipv6 addresses
func IpToUint(realIp string) uint32 {
	netIp := net.ParseIP(realIp)
	if netIp == nil {
		return 0
	}

	var ip uint32
	netIp.To4()
	if len(netIp) == net.IPv6len {
		ip = binary.BigEndian.Uint32(netIp[12:16])
	} else if len(netIp) == net.IPv4len {
		ip = binary.BigEndian.Uint32(netIp)
	}
	return ip
}

func Ip4ToString(ip uint32) string {
	netIp := make(net.IP, 4)
	binary.BigEndian.PutUint32(netIp, ip)
	return netIp.String()
}

func GetIp4Block(ip uint32) uint32 {
	return ip & 0xFFFFFF00
}

func GetIp4BlockStr(ip string) string {
	ipBlock := GetIp4Block(IpToUint(ip))
	return Ip4ToString(ipBlock)
}

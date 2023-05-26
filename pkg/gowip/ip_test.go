package gowip

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIpv4Int(t *testing.T) {
	ip := IpToUint("127.0.0.1")
	assert.Equal(t, uint32(0x7f000001), ip)
}

func TestIpv6Int(t *testing.T) {
	ip := IpToUint("0000:0000:0000:0000:0000:0000:0000:0001")
	assert.Equal(t, uint32(0x1), ip)
}

func TestIpToUintWithInvalidInput(t *testing.T) {
	ip := IpToUint("!3s nope this is not an ip")
	assert.Equal(t, uint32(0), ip)
}

func TestGetIp4Block(t *testing.T) {
	ip := IpToUint("127.0.0.1")
	assert.Equal(t, uint32(0x7f000001), ip)

	ipBlock := GetIp4Block(ip)
	assert.Equal(t, uint32(0x7f000000), ipBlock)
}

func TestIpToString(t *testing.T) {
	ip := Ip4ToString(0x7f000001)
	assert.Equal(t, "127.0.0.1", ip)

	ip = Ip4ToString(0x00000000)
	assert.Equal(t, "0.0.0.0", ip)

	ip = Ip4ToString(0xFFFFFFFF)
	assert.Equal(t, "255.255.255.255", ip)
}

func TestGetIp4BlockStr(t *testing.T) {
	ipBlock := GetIp4BlockStr("127.0.0.1")
	assert.Equal(t, "127.0.0.0", ipBlock)

	ipBlock = GetIp4BlockStr("255.255.255.255")
	assert.Equal(t, "255.255.255.0", ipBlock)
}

package webutil

import (
	"net"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestLocalIP(t *testing.T) {
	assert := assert.New(t)

	// Without any patching of default.
	assert.NotEmpty(LocalIP())

	defer restoreGetInterfaceAddrs()
	_, localIP, err := net.ParseCIDR("192.168.7.31/32")
	assert.Nil(err)
	_, loopbackIP, err := net.ParseCIDR("127.0.0.1/32")
	assert.Nil(err)
	_, ipv6, err := net.ParseCIDR("2001:db8::/32")
	assert.Nil(err)

	// Happy path: Local IP returned (i.e. non-loopback)
	setGetInterfaceAddrs(func() ([]net.Addr, error) {
		return []net.Addr{ipv6, localIP, loopbackIP}, nil
	})
	assert.Equal("192.168.7.31", LocalIP())

	// Loopback IP returned
	setGetInterfaceAddrs(func() ([]net.Addr, error) {
		return []net.Addr{ipv6, loopbackIP}, nil
	})
	assert.Equal("127.0.0.1", LocalIP())

	// No interface addresses
	setGetInterfaceAddrs(func() ([]net.Addr, error) {
		return nil, ex.New("big failure")
	})
	assert.Equal("", LocalIP())

	// No interface addresses
	setGetInterfaceAddrs(func() ([]net.Addr, error) {
		return nil, nil
	})
	assert.Equal("", LocalIP())
}

package statsd

import (
	"net"
	"strings"
)

// NewUDPListener returns a new UDP listener for a given address.
func NewUDPListener(addr string) (net.PacketConn, error) {
	listener, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

// NewUnixgramListener returns a new unixgram listener for a given path.
func NewUnixgramListener(path string) (net.PacketConn, error) {
	path = strings.TrimPrefix(path, "unix://")
	listener, err := net.ListenPacket("unixgram", path)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

package statsd

import "net"

// NewUDPListener returns a new UDP listener for a given address.
func NewUDPListener(addr string) (net.PacketConn, error) {
	listener, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

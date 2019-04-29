package webutil

import (
	"net"
	"time"
)

var (
	_ net.Listener = (*TCPKeepAliveListener)(nil)
)

// Defaults
const (
	DefaultTCPKeepAlivePeriod = 3 * time.Minute
)

// TCPKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
// Taken from net/http/server.go
type TCPKeepAliveListener struct {
	*net.TCPListener

	KeepAlivePeriod time.Duration
}

// TCPKeepAlivePeriodOrDefault returns the KeepAlivePeriod or a default.
func (ln TCPKeepAliveListener) TCPKeepAlivePeriodOrDefault() time.Duration {
	if ln.KeepAlivePeriod > 0 {
		return ln.KeepAlivePeriod
	}
	return DefaultTCPKeepAlivePeriod
}

// Accept implements net.Listener
func (ln TCPKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(ln.TCPKeepAlivePeriodOrDefault())
	return tc, nil
}

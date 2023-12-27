/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webutil

import (
	"net"
	"time"
)

var (
	_ net.Listener = (*TCPKeepAliveListener)(nil)
)

// TCPKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
//
// It is taken from net/http/server.go
type TCPKeepAliveListener struct {
	*net.TCPListener

	KeepAlive	bool
	KeepAlivePeriod	time.Duration
}

// Accept implements net.Listener
func (ln TCPKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	_ = tc.SetKeepAlive(ln.KeepAlive)
	_ = tc.SetKeepAlivePeriod(ln.KeepAlivePeriod)
	return tc, nil
}

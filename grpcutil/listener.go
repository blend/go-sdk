package grpcutil

import (
	"net"
	"strings"
	"github.com/blend/go-sdk/ex"
)

// Listener creates a net listener for a given bind address.
// It handles detecting if we should create a unix socket address.
func Listener(bindAddr string) (net.Listener, error) {
	var socketListener net.Listener
	var err error
	if strings.HasPrefix(bindAddr, "unix://") {
		socketListener, err = net.Listen("unix", strings.TrimPrefix(bindAddr, "unix://"))
	} else {
		socketListener, err = net.Listen("tcp", bindAddr)
	}
	return socketListener, ex.New(err)
}

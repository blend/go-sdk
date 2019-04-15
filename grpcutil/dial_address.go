package grpcutil

import (
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"

	"github.com/blend/go-sdk/ex"
)

// DialAddress dials an address with a given set of dial options.
// It resolves how to dial unix sockets if the address is prefixed with `unix://`.
func DialAddress(addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	if strings.HasPrefix("unix://", addr) {
		opts = append(opts,
			grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
				return net.DialTimeout("unix", addr, timeout)
			}))
		addr = strings.TrimPrefix(addr, "unix://")
	}

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, ex.New(err)
	}
	return conn, nil
}

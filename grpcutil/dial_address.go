package grpcutil

import (
	"context"
	"net"
	"strings"

	"google.golang.org/grpc"

	"github.com/blend/go-sdk/ex"
)

// DialAddress dials an address with a given set of dial options.
// It resolves how to dial unix sockets if the address is prefixed with `unix://`.
func DialAddress(addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	if strings.HasPrefix("unix://", addr) {
		opts = append(opts,
			grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
				return new(net.Dialer).DialContext(ctx, "unix", addr)
			}))
		addr = strings.TrimPrefix(addr, "unix://")
	}

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, ex.New(err)
	}
	return conn, nil
}

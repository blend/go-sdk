package proxyprotocol

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/blend/go-sdk/webutil"
)

// CreateListener creates a new proxy protocol listener.
func CreateListener(addr string, opts ...CreateListenerOption) (net.Listener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	options := CreateListenerOptions{
		KeepAlive:       true,
		KeepAlivePeriod: 3 * time.Minute,
	}
	for _, opt := range opts {
		if err := opt(&options); err != nil {
			return nil, err
		}
	}

	var output net.Listener = webutil.TCPKeepAliveListener{TCPListener: ln.(*net.TCPListener)}

	if options.TLSConfig != nil {
		output = tls.NewListener(output, options.TLSConfig)
	}
	if options.UseProxyProtocol {
		output = &Listener{Listener: output}
	}
	return output, nil
}

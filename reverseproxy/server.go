package reverseproxy

import (
	"crypto/tls"
	"net"
	"net/http"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/proxyprotocol"
	"github.com/blend/go-sdk/webutil"
)

// NewServer returns a new server.
func NewServer(opts ...ServerOption) (*Server, error) {
	svr := Server{
		Latch: async.NewLatch(),
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				// infosec approved cipher suites
				// order matters here, DO NOT MOVE the first cipher lower it is required
				// for http 2 that this be the first ciper in the list
				// https://github.com/golang/go/issues/20213
				// ciphers are dark magic and chrome is mean
				// https://support.cloudflare.com/hc/en-us/articles/200933580-What-cipher-suites-does-CloudFlare-use-for-SSL-
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			PreferServerCipherSuites: true},
		Server: &http.Server{},
	}

	var err error
	for _, opt := range opts {
		if err = opt(&svr); err != nil {
			return nil, err
		}
	}
	return &svr, nil
}

// ServerOption is a mutator for the server.
type ServerOption func(*Server) error

// Server represents an opinionated http server host for the reverse proxy.
type Server struct {
	*async.Latch
	Proxy

	TLSConfig      *tls.Config
	Server         *http.Server
	RedirectServer *http.Server
}

// CreateListener creates a new proxy protocol listener.
func (s Server) CreateListener(addr string) (*proxyprotocol.Listener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &proxyprotocol.Listener{Listener: webutil.TCPKeepAliveListener{TCPListener: ln.(*net.TCPListener)}}, nil
}

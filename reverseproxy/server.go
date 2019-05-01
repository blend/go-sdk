package reverseproxy

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/certutil"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/proxyprotocol"
	"github.com/blend/go-sdk/webutil"
)

// NewServer returns a new server.
func NewServer(proxy *Proxy, opts ...ServerOption) (*Server, error) {
	svr := Server{
		Latch: async.NewLatch(),
		Proxy: proxy,
	}

	var err error
	for _, opt := range opts {
		if err = opt(&svr); err != nil {
			return nil, err
		}
	}

	return &svr, nil
}

// Server represents an opinionated http server host for the reverse proxy.
// It serves TLS traffic, upgrades 80->443, and watches for cert changes.
type Server struct {
	*async.Latch

	UseProxyProtocol bool

	DefaultMiddleware []webutil.Middleware
	Log               logger.Log
	Proxy             *Proxy
	CertWatcher       *certutil.CertWatcher
	Server            *http.Server
	UpgradeServer     *http.Server
}

// CreateListener creates a new proxy protocol listener.
func (s Server) CreateListener(addr string) (net.Listener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	var output net.Listener = webutil.TCPKeepAliveListener{TCPListener: ln.(*net.TCPListener)}
	if s.Server.TLSConfig != nil {
		output = tls.NewListener(output, s.Server.TLSConfig)
	}
	if s.UseProxyProtocol {
		output = &proxyprotocol.Listener{Listener: output}
	}
	return output, nil
}

// Upgrade is the http handler to upgrade http request to https.
func (s Server) Upgrade(rw http.ResponseWriter, req *http.Request) {
	w := NewResponseWriter(rw)

	if s.Log != nil {
		s.Log.Trigger(req.Context(), logger.NewHTTPRequestEvent(req))

		start := time.Now()
		defer func() {
			wre := logger.NewHTTPResponseEvent(req,
				logger.OptHTTPResponseStatusCode(http.StatusMovedPermanently),
				logger.OptHTTPResponseContentLength(w.ContentLength()),
				logger.OptHTTPResponseElapsed(time.Since(start)),
			)

			if value := rw.Header().Get("Content-Type"); len(value) > 0 {
				wre.ContentType = value
			}
			if value := rw.Header().Get("Content-Encoding"); len(value) > 0 {
				wre.ContentEncoding = value
			}
			s.Log.Trigger(req.Context(), wre)
		}()
	}

	req.URL.Scheme = schemeHTTPS
	if req.URL.Host == "" {
		req.URL.Host = req.Host
	}

	http.Redirect(rw, req, req.URL.String(), http.StatusMovedPermanently)
}

// Start starts the server.
func (s *Server) Start() error {
	s.Starting()
	var err error

	if s.Server == nil {
		s.Server = &http.Server{
			Addr: DefaultAddr,
		}
	}

	if s.Proxy.Log == nil {
		s.Proxy.Log = s.Log
	}
	for _, upstream := range s.Proxy.Upstreams {
		transport := &http.Transport{}
		if err = http2.ConfigureTransport(transport); err != nil {
			return err
		}
		upstream.ReverseProxy.Transport = transport
		if upstream.Log == nil {
			upstream.Log = s.Log
		}
	}
	s.Server.Handler = http.HandlerFunc(webutil.NestMiddleware(s.Proxy.ServeHTTP, s.DefaultMiddleware...))

	serverListener, err := s.CreateListener(s.Server.Addr)
	if err != nil {
		return err
	}

	actions := []func() error{
		func() error {
			s.Started()
			return s.Server.Serve(serverListener)
		},
	}

	if s.CertWatcher != nil {
		actions = append(actions, s.CertWatcher.Start)
	}
	if s.UpgradeServer != nil {
		actions = append(actions, s.UpgradeServer.ListenAndServe)
	}

	if err = async.RunToError(actions...); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop stops the server
func (s *Server) Stop() error {
	if !s.CanStop() {
		return async.ErrCannotStop
	}

	if s.CertWatcher != nil {
		if err := s.CertWatcher.Stop(); err != nil {
			return err
		}
		logger.MaybeInfo(s.Log, "cert watcher stopped")
	}

	if s.UpgradeServer != nil {
		deadline, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.UpgradeServer.Shutdown(deadline); err != nil && err != context.Canceled {
			return err
		}
		logger.MaybeInfo(s.Log, "upgrade server stopped")
	}

	if s.Server != nil {
		deadline, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Server.Shutdown(deadline); err != nil && err != context.Canceled {
			return err
		}
		logger.MaybeInfo(s.Log, "reverse proxy server stopped")
	}

	logger.MaybeInfo(s.Log, "reverse proxy stopped")
	return nil
}

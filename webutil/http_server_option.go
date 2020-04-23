package webutil

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"
)

// HTTPServerOption is a mutator for an http server.
type HTTPServerOption func(*http.Server) error

// OptHTTPServerHandler mutates a http server.
func OptHTTPServerHandler(handler http.Handler) HTTPServerOption {
	return func(s *http.Server) error {
		s.Handler = handler
		return nil
	}
}

// OptHTTPServerTLSConfig mutates a http server.
func OptHTTPServerTLSConfig(cfg *tls.Config) HTTPServerOption {
	return func(s *http.Server) error {
		s.TLSConfig = cfg
		return nil
	}
}

// OptHTTPServerAddr mutates a http server.
func OptHTTPServerAddr(addr string) HTTPServerOption {
	return func(s *http.Server) error {
		s.Addr = addr
		return nil
	}
}

// OptHTTPServerMaxHeaderBytes mutates a http server.
func OptHTTPServerMaxHeaderBytes(value int) HTTPServerOption {
	return func(s *http.Server) error {
		s.MaxHeaderBytes = value
		return nil
	}
}

// OptHTTPServerReadTimeout mutates a http server.
func OptHTTPServerReadTimeout(value time.Duration) HTTPServerOption {
	return func(s *http.Server) error {
		s.ReadTimeout = value
		return nil
	}
}

// OptHTTPServerReadHeaderTimeout mutates a http server.
func OptHTTPServerReadHeaderTimeout(value time.Duration) HTTPServerOption {
	return func(s *http.Server) error {
		s.ReadHeaderTimeout = value
		return nil
	}
}

// OptHTTPServerWriteTimeout mutates a http server.
func OptHTTPServerWriteTimeout(value time.Duration) HTTPServerOption {
	return func(s *http.Server) error {
		s.WriteTimeout = value
		return nil
	}
}

// OptHTTPServerIdleTimeout mutates a http server.
func OptHTTPServerIdleTimeout(value time.Duration) HTTPServerOption {
	return func(s *http.Server) error {
		s.IdleTimeout = value
		return nil
	}
}

// OptHTTPServerErrorLog sets the error log.
func OptHTTPServerErrorLog(log *log.Logger) HTTPServerOption {
	return func(s *http.Server) error {
		s.ErrorLog = log
		return nil
	}
}

/*
// WCTODO: once we bump to 1.13
// OptHTTPServerBaseContext sets the error log.
func OptHTTPServerBaseContext(handler func(net.Listener) context.Context) HTTPServerOption {
	return func(s *http.Server) error {
		s.BaseContext = handler
		return nil
	}
}
*/

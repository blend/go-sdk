/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webutil

import (
	"context"
	"crypto/tls"
	"log"
	"net"
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

// OptHTTPServerBaseContext sets the base context for requests to a given server.
func OptHTTPServerBaseContext(baseContextProvider func(net.Listener) context.Context) HTTPServerOption {
	return func(s *http.Server) error {
		s.BaseContext = baseContextProvider
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

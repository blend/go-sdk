package reverseproxy

import (
	"crypto/tls"
	"net/http"

	"github.com/blend/go-sdk/certutil"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

// ServerOption is a mutator for the server.
type ServerOption func(*Server) error

// OptServerLog sets the server logger.
func OptServerLog(log logger.Log) ServerOption {
	return func(s *Server) error {
		s.Log = log
		return nil
	}
}

// OptServerUseProxyProtocol sets if we should use a proxy protocol listener.
func OptServerUseProxyProtocol(useProxyProtocol bool) ServerOption {
	return func(s *Server) error {
		s.UseProxyProtocol = useProxyProtocol
		return nil
	}
}

// OptServerMiddleware sets the server default middleware.
func OptServerMiddleware(middleware ...webutil.Middleware) ServerOption {
	return func(s *Server) error {
		s.DefaultMiddleware = middleware
		return nil
	}
}

// OptServerTLS sets the server to listen with tls.
func OptServerTLS(certPath, keyPath string) ServerOption {
	return func(s *Server) error {
		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			return err
		}

		if s.Server == nil {
			s.Server = &http.Server{}
		}

		if s.Server.TLSConfig == nil {
			s.Server.TLSConfig = &tls.Config{}
		}
		s.Server.TLSConfig.Certificates = append(s.Server.TLSConfig.Certificates, cert)
		return nil
	}
}

// OptServerTLSCertWatcher sets up the cert watcher for the server.
func OptServerTLSCertWatcher(certPath, keyPath string) ServerOption {
	return func(s *Server) error {
		var err error
		s.CertWatcher, err = certutil.NewCertWatcher(certPath, keyPath)
		if err != nil {
			return err
		}

		if s.Server == nil {
			s.Server = &http.Server{}
		}

		if s.Server.TLSConfig == nil {
			s.Server.TLSConfig = &tls.Config{}
		}

		s.Server.TLSConfig.GetCertificate = s.CertWatcher.GetCertificate
		return nil
	}
}

// OptServerBindAddr sets the server bind address.
func OptServerBindAddr(bindAddr string) ServerOption {
	return func(s *Server) error {
		if s.Server == nil {
			s.Server = &http.Server{}
		}
		s.Server.Addr = bindAddr
		return nil
	}
}

// OptServerUpgrade sets the server upgrade bind address and redirects http requests to https.
func OptServerUpgrade(bindAddr string) ServerOption {
	return func(s *Server) error {
		if s.UpgradeServer == nil {
			s.UpgradeServer = &http.Server{}
		}
		s.UpgradeServer.Addr = bindAddr
		s.UpgradeServer.Handler = http.HandlerFunc(s.Upgrade)
		return nil
	}
}

// OptServerTLSCipherSuites sets the server to use secure cipher suites.
func OptServerTLSCipherSuites() ServerOption {
	return func(s *Server) error {
		if s.Server == nil {
			s.Server = &http.Server{}
		}
		if s.Server.TLSConfig == nil {
			s.Server.TLSConfig = &tls.Config{}
		}
		s.Server.TLSConfig.MinVersion = tls.VersionTLS12
		s.Server.TLSConfig.CipherSuites = []uint16{
			// Order matters here, DO NOT MOVE the first cipher lower it is required
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
		}
		s.Server.TLSConfig.PreferServerCipherSuites = true
		return nil
	}
}

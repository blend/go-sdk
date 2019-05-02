package webutil

import "crypto/tls"

// TLSSecureCipherSuites sets the tls config to use secure cipher suites.
func TLSSecureCipherSuites(tlsConfig *tls.Config) {
	tlsConfig.MinVersion = tls.VersionTLS12
	tlsConfig.CipherSuites = []uint16{
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
	tlsConfig.PreferServerCipherSuites = true
}

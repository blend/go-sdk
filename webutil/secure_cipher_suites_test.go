package webutil

import (
	"crypto/tls"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestTLSSecureCipherSuites(t *testing.T) {
	assert := assert.New(t)

	c := &tls.Config{}
	// Before
	assert.Equal(0, c.MinVersion)
	assert.Nil(c.CipherSuites)
	assert.False(c.PreferServerCipherSuites)

	TLSSecureCipherSuites(c)
	// After
	assert.Equal(tls.VersionTLS12, c.MinVersion)
	expected := []uint16{
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
	assert.Equal(expected, c.CipherSuites)
	assert.True(c.PreferServerCipherSuites)
}

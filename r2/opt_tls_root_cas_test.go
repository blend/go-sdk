package r2

import (
	"crypto/x509"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptTLSRootCAs(t *testing.T) {
	assert := assert.New(t)

	pool, err := x509.SystemCertPool()
	assert.Nil(err)
	r := New("http://foo.com", OptTLSRootCAs(pool))
	assert.NotNil(r.Client.Transport.(*http.Transport).TLSClientConfig.RootCAs)
}

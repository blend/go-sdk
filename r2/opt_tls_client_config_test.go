package r2

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptTLSCLientConfig(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptTLSClientConfig(&tls.Config{}))
	assert.NotNil(r.Client.Transport.(*http.Transport).TLSClientConfig)
}

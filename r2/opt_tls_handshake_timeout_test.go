package r2

import (
	"net/http"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestOptTLSHandshakeTimeout(t *testing.T) {
	assert := assert.New(t)

	r := New("http://foo.com", OptTLSHandshakeTimeout(time.Second))
	assert.Equal(time.Second, r.Client.Transport.(*http.Transport).TLSHandshakeTimeout)
}

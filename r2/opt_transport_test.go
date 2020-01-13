package r2

import (
	"net/http"
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestOptTransport(t *testing.T) {
	assert := assert.New(t)

	r := New("http://foo.com", OptTransport(&http.Transport{}))
	assert.NotNil(r.Client.Timeout)
}

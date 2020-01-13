package r2

import (
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestOptPort(t *testing.T) {
	assert := assert.New(t)

	r := New("http://foo.com", OptPort(8443))
	assert.Equal("http://foo.com:8443", r.URL.String())
}

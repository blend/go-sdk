package r2

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptScheme(t *testing.T) {
	assert := assert.New(t)

	r := New("http://foo.com", OptScheme("spdy"))
	assert.Equal("spdy", r.Request.URL.Scheme)

	var unset Request
	OptScheme("spdy")(&unset)
	assert.Equal("spdy", unset.Request.URL.Scheme)
}

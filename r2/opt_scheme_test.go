package r2

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptScheme(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptScheme("spdy"))
	assert.Equal("spdy", r.Request.URL.Scheme)

	var unset Request
	assert.NotNil(OptScheme("spdy")(&unset))
}

package r2

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptHost(t *testing.T) {
	assert := assert.New(t)

	r := New("http://foo.com", OptHost("bar.com"))

	assert.Equal("http://bar.com", r.URL.String())
}

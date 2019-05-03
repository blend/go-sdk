package r2

import (
	"strings"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptBasicAuth(t *testing.T) {
	assert := assert.New(t)

	opt := OptBasicAuth("foo", "bar")

	req := New("https://foo.bar.local")
	assert.Nil(opt(req))

	assert.NotNil(req.Header)
	assert.NotEmpty(req.Header.Get("Authorization"))
	assert.True(strings.HasPrefix(req.Header.Get("Authorization"), "Basic "))
}

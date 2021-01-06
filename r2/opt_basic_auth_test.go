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

	assert.NotNil(req.Request.Header)
	assert.NotEmpty(req.Request.Header.Get("Authorization"))
	assert.True(strings.HasPrefix(req.Request.Header.Get("Authorization"), "Basic "))
}

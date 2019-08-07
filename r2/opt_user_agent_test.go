package r2

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptUserAgent(t *testing.T) {
	assert := assert.New(t)

	opt := OptUserAgent("blend test harness")
	req := New("http://foo.bar.local")
	assert.NotEqual("blend test harness", req.UserAgent())
	assert.Nil(opt(req))
	assert.Equal("blend test harness", req.UserAgent())
}

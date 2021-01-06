package r2

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptHost(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptHost("bar.invalid"))

	assert.Equal("https://bar.invalid/test?query=value", r.Request.URL.String())
}

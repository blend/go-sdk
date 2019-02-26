package r2

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptMethods(t *testing.T) {
	assert := assert.New(t)

	req := New("https://foo.bar.local")

	OptMethod("OPTIONS")(req)
	assert.Equal("OPTIONS", req.Method)

	OptGet()(req)
	assert.Equal("GET", req.Method)

	OptPost()(req)
	assert.Equal("POST", req.Method)

	OptPut()(req)
	assert.Equal("PUT", req.Method)

	OptPatch()(req)
	assert.Equal("PATCH", req.Method)

	OptDelete()(req)
	assert.Equal("DELETE", req.Method)
}

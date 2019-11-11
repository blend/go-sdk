package webutil

import (
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestGetPort(t *testing.T) {
	assert := assert.New(t)

	assert.Empty(GetPort(nil))
	assert.Empty(GetPort(&http.Request{}))
	assert.Equal("8443", GetPort(&http.Request{
		Header: http.Header{
			HeaderXForwardedPort: {"8443"},
		},
	}), "should use existing header if found")
	assert.Equal("8443", GetPort(&http.Request{
		Header: http.Header{
			HeaderXForwardedPort: {"9090,8443"},
		},
	}), "should use existing header last chunk if found")
}

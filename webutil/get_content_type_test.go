package webutil

import (
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestGetContentType(t *testing.T) {
	assert := assert.New(t)

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	assert.Equal("application/json", GetContentType(headers))

	headers = http.Header{}
	headers.Set("X-Forwarded-Host", "local.bar.com")
	assert.Equal("", GetContentType(headers))

	assert.Equal("", GetContentType(nil))
}

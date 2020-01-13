package webutil

import (
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestGetContentEncoding(t *testing.T) {
	assert := assert.New(t)

	headers := http.Header{}
	headers.Set("Content-Encoding", "gzip")
	assert.Equal("gzip", GetContentEncoding(headers))

	headers = http.Header{}
	headers.Set("Content-Type", "application/json")
	assert.Equal("", GetContentEncoding(headers))

	assert.Equal("", GetContentEncoding(nil))
}

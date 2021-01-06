package web

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestSession(t *testing.T) {
	assert := assert.New(t)

	session := &Session{}
	session.WithBaseURL("https://foo.com/bar")
	assert.Equal("https://foo.com/bar", session.BaseURL)
	session.WithUserAgent("example-string")
	assert.Equal("example-string", session.UserAgent)
	session.WithRemoteAddr("10.10.32.1")
	assert.Equal("10.10.32.1", session.RemoteAddr)
}

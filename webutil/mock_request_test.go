package webutil

import (
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestNewMockRequest(t *testing.T) {
	assert := assert.New(t)

	req := NewMockRequest("OPTIONS", "/foo")
	assert.Equal("OPTIONS", req.Method)
	assert.Equal("localhost:8080", req.Host)
	assert.Equal("/foo", req.RequestURI)
	assert.NotNil(req.URL)
	assert.Equal("http", req.Proto)
	assert.Equal(1, req.ProtoMajor)
	assert.Equal(1, req.ProtoMinor)
	assert.Equal("http", req.URL.Scheme)
	assert.Equal("/foo", req.URL.Path)
	assert.Equal("127.0.0.1:8080", req.RemoteAddr)
	assert.NotNil(req.Header)
	assert.Equal("go-sdk test", req.Header.Get(HeaderUserAgent))
}

func TestNewMockRequestWithCookie(t *testing.T) {
	assert := assert.New(t)
	req := NewMockRequestWithCookie("OPTIONS", "/foo", "foo", "bar")
	assert.NotEmpty(req.Cookies())

	cookie, err := req.Cookie("foo")
	assert.Nil(err)
	assert.Equal("foo", cookie.Name)
	assert.Equal("bar", cookie.Value)
}

package webutil

import (
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestGetRemoteAddr(t *testing.T) {
	assert := assert.New(t)

	hdr := http.Header{}
	hdr.Set("X-Forwarded-For", "123")
	r := http.Request{
		Header: hdr,
	}
	assert.Equal("123", GetRemoteAddr(&r))

	hdr = http.Header{}
	hdr.Set("X-FORWARDED-FOR", "1")
	r = http.Request{
		Header: hdr,
	}
	assert.Equal("1", GetRemoteAddr(&r))

	hdr = http.Header{}
	hdr.Set("X-FORWARDED-FOR", "1,2,3")
	r = http.Request{
		Header: hdr,
	}
	assert.Equal("3", GetRemoteAddr(&r))

	hdr = http.Header{}
	hdr.Set("X-Real-Ip", "1")
	r = http.Request{
		Header: hdr,
	}
	assert.Equal("1", GetRemoteAddr(&r))

	hdr = http.Header{}
	hdr.Set("X-REAL-IP", "1")
	r = http.Request{
		Header: hdr,
	}
	assert.Equal("1", GetRemoteAddr(&r))

	hdr = http.Header{}
	hdr.Set("X-REAL-IP", "1,2,3")
	r = http.Request{
		Header: hdr,
	}
	assert.Equal("3", GetRemoteAddr(&r))

	r = http.Request{
		RemoteAddr: "1:1",
	}
	assert.Equal("1", GetRemoteAddr(&r))

	r = http.Request{
		RemoteAddr: "1",
	}
	assert.Equal("", GetRemoteAddr(&r))
}

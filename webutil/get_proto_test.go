package webutil

import (
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestGetProto(t *testing.T) {
	assert := assert.New(t)

	headers := http.Header{}
	headers.Set(HeaderXForwardedProto, SchemeHTTPS)
	r := http.Request{
		Proto:  SchemeHTTP + "/1.0",
		Header: headers,
	}
	assert.Equal(SchemeHTTPS, GetProto(&r))

	headers = http.Header{}
	headers.Set(HeaderXForwardedProto, SchemeSPDY+","+SchemeHTTPS)
	r = http.Request{
		Proto:  SchemeHTTP + "/1.0",
		Header: headers,
	}
	assert.Equal(SchemeHTTPS, GetProto(&r))

	headers = http.Header{}
	r = http.Request{
		Header: headers,
	}
	assert.Empty(GetProto(&r))
}

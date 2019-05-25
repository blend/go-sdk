package logger

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestWriteHTTPRequest(t *testing.T) {
	assert := assert.New(t)

	tf := TextOutputFormatter{
		NoColor: true,
	}
	buf := new(bytes.Buffer)
	WriteHTTPRequest(tf, buf, &http.Request{Method: "GET", URL: &url.URL{Scheme: "http", Host: "localhost", Path: "/foo"}})
	assert.Equal("GET /foo", buf.String())
}

func TestWriteHTTPResponse(t *testing.T) {
	assert := assert.New(t)

	tf := TextOutputFormatter{
		NoColor: true,
	}
	buf := new(bytes.Buffer)
	req := &http.Request{Method: "GET", URL: &url.URL{Scheme: "http", Host: "localhost", Path: "/foo"}}
	WriteHTTPResponse(tf, buf, req, http.StatusOK, 1024, "text/html", time.Second)
	assert.Equal("GET http://localhost/foo 200 1s text/html 1kB", buf.String())
}

func TestWriteFields(t *testing.T) {
	assert := assert.New(t)

	tf := TextOutputFormatter{
		NoColor: true,
	}
	buf := new(bytes.Buffer)

	WriteFields(tf, buf, Fields{"foo": "bar", "moo": "loo"})
	assert.Contains(buf.String(), "foo=bar")
	assert.Contains(buf.String(), "moo=loo")
}

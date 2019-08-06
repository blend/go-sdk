package logger

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/blend/go-sdk/ansi"
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

func TestWriteHTTPRequestPrefersRawPath(t *testing.T) {
	assert := assert.New(t)

	tf := TextOutputFormatter{
		NoColor: true,
	}
	buf := new(bytes.Buffer)
	WriteHTTPRequest(tf, buf, &http.Request{Method: "GET", URL: &url.URL{Scheme: "http", Host: "localhost", RawPath: "/foo?hi=there", Path: "/foo"}})
	assert.Equal("GET /foo?hi=there", buf.String())
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

func TestFormatFields(t *testing.T) {
	assert := assert.New(t)

	tf := NewTextOutputFormatter(OptTextNoColor())
	actual := FormatFields(tf, ansi.ColorBlue, Fields{"foo": "bar", "moo": "loo"})
	assert.Equal("foo=bar moo=loo", actual)

	actual = FormatFields(tf, ansi.ColorBlue, Fields{"moo": "loo", "foo": "bar"})
	assert.Equal("foo=bar moo=loo", actual)

	tf = NewTextOutputFormatter()
	actual = FormatFields(tf, ansi.ColorBlue, Fields{"foo": "bar", "moo": "loo"})
	assert.Equal(ansi.ColorBlue.Apply("foo")+"=bar "+ansi.ColorBlue.Apply("moo")+"=loo", actual)

	actual = FormatFields(tf, ansi.ColorBlue, Fields{"moo": "loo", "foo": "bar"})
	assert.Equal(ansi.ColorBlue.Apply("foo")+"=bar "+ansi.ColorBlue.Apply("moo")+"=loo", actual)
}

func TestFormatHeaders(t *testing.T) {
	assert := assert.New(t)

	tf := NewTextOutputFormatter(OptTextNoColor())
	actual := FormatHeaders(tf, ansi.ColorBlue, http.Header{"Foo": []string{"bar"}, "Moo": []string{"loo"}})
	assert.Equal("{ Foo:bar Moo:loo }", actual)

	actual = FormatHeaders(tf, ansi.ColorBlue, http.Header{"Moo": []string{"loo"}, "Foo": []string{"bar"}})
	assert.Equal("{ Foo:bar Moo:loo }", actual)

	tf = NewTextOutputFormatter()
	actual = FormatHeaders(tf, ansi.ColorBlue, http.Header{"Foo": []string{"bar"}, "Moo": []string{"loo"}})
	assert.Equal("{ "+ansi.ColorBlue.Apply("Foo")+":bar "+ansi.ColorBlue.Apply("Moo")+":loo }", actual)
}

package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/assert"
)

func TestNewHTTPResponseEvent(t *testing.T) {
	assert := assert.New(t)

	hre := NewHTTPResponseEvent(nil,
		OptHTTPResponseMeta(OptEventMetaFlagColor(ansi.ColorBlue)),
		OptHTTPResponseRequest(&http.Request{Method: "foo", URL: &url.URL{Scheme: "https", Host: "localhost", Path: "/foo/bailey"}}),
		OptHTTPResponseContentEncoding("utf-8"),
		OptHTTPResponseContentLength(1337),
		OptHTTPResponseContentType("text/html"),
		OptHTTPResponseElapsed(time.Second),
		OptHTTPResponseRoute("/foo/:bar"),
		OptHTTPResponseStatusCode(http.StatusOK),
		OptHTTPResponseState("this is the state"),
	)

	assert.Equal(ansi.ColorBlue, hre.GetFlagColor())
	assert.Equal("foo", hre.Request.Method)
	assert.Equal("utf-8", hre.ContentEncoding)
	assert.Equal(1337, hre.ContentLength)
	assert.Equal("text/html", hre.ContentType)
	assert.Equal(time.Second, hre.Elapsed)
	assert.Equal("/foo/:bar", hre.Route)
	assert.Equal(http.StatusOK, hre.StatusCode)
	assert.Equal("this is the state", hre.State)

	noColor := NewTextOutputFormatter(OptTextNoColor())
	buf := new(bytes.Buffer)
	hre.WriteText(noColor, buf)
	assert.NotContains(buf.String(), "/foo/:bar")
	assert.Contains(buf.String(), "/foo/bailey")

	contents, err := json.Marshal(hre)
	assert.Nil(err)
	assert.Contains(string(contents), "/foo/:bar")
}

func TestHTTPResponseEventListener(t *testing.T) {
	assert := assert.New(t)

	var didCall bool
	listener := NewHTTPResponseEventListener(func(_ context.Context, hre *HTTPResponseEvent) {
		didCall = true
	})
	listener(context.Background(), NewMessageEvent(Info, "test"))
	assert.False(didCall)
	listener(context.Background(), NewHTTPResponseEvent(nil))
	assert.True(didCall)
}

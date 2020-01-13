package webutil

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"
	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
)

func TestNewHTTPResponseEvent(t *testing.T) {
	assert := assert.New(t)

	hre := NewHTTPResponseEvent(nil,
		OptHTTPResponseRequest(&http.Request{Method: "foo", URL: &url.URL{Scheme: "https", Host: "localhost", Path: "/foo/bailey"}}),
		OptHTTPResponseContentEncoding("utf-8"),
		OptHTTPResponseContentLength(1337),
		OptHTTPResponseContentType("text/html"),
		OptHTTPResponseElapsed(time.Second),
		OptHTTPResponseRoute("/foo/:bar"),
		OptHTTPResponseStatusCode(http.StatusOK),
		OptHTTPResponseHeader(http.Header{"X-Bad": []string{"nope", "definitely nope"}}),
	)

	assert.Equal("foo", hre.Request.Method)
	assert.Equal("utf-8", hre.ContentEncoding)
	assert.Equal(1337, hre.ContentLength)
	assert.Equal("text/html", hre.ContentType)
	assert.Equal(time.Second, hre.Elapsed)
	assert.Equal("/foo/:bar", hre.Route)
	assert.Equal(http.StatusOK, hre.StatusCode)
	assert.Equal("nope", hre.Header.Get("X-Bad"))

	noColor := logger.NewTextOutputFormatter(logger.OptTextNoColor())
	buf := new(bytes.Buffer)
	hre.WriteText(noColor, buf)
	assert.NotContains(buf.String(), "/foo/:bar")
	assert.Contains(buf.String(), "/foo/bailey")
	assert.NotContains(buf.String(), "X-Bad", "response headers should not be written to text output")
	assert.NotContains(buf.String(), "definitely nope", "response headers should not be written to text output")

	contents, err := json.Marshal(hre.Decompose())
	assert.Nil(err)
	assert.Contains(string(contents), "/foo/:bar")

	assert.NotContains(string(contents), "X-Bad", "response headers should not be written to json output")
	assert.NotContains(string(contents), "definitely nope", "response headers should not be written to json output")
}

func TestHTTPResponseEventListener(t *testing.T) {
	assert := assert.New(t)

	var didCall bool
	listener := NewHTTPResponseEventListener(func(_ context.Context, hre HTTPResponseEvent) {
		didCall = true
	})
	listener(context.Background(), logger.NewMessageEvent(logger.Info, "test"))
	assert.False(didCall)
	listener(context.Background(), NewHTTPResponseEvent(nil))
	assert.True(didCall)
}

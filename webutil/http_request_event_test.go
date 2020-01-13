package webutil

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
)

func TestNewHTTPRequestEvent(t *testing.T) {
	assert := assert.New(t)

	hre := NewHTTPRequestEvent(nil,
		OptHTTPRequest(&http.Request{
			Method: "GET",
			URL: &url.URL{
				Scheme: "https",
				Host:   "localhost",
				Path:   "/foo",
			},
		}),
		OptHTTPRequestRoute("/foo/:bar"),
	)

	assert.NotNil(hre.Request)
	assert.Equal("GET", hre.Request.Method)
	assert.Equal("/foo", hre.Request.URL.Path)
	assert.Equal("/foo/:bar", hre.Route)

	noColor := logger.TextOutputFormatter{
		NoColor: true,
	}

	buf := new(bytes.Buffer)
	hre.WriteText(noColor, buf)
	assert.NotEmpty(buf.String())

	contents, err := json.Marshal(hre.Decompose())
	assert.Nil(err)
	assert.NotEmpty(contents)
}

func TestHTTPRequestEventListener(t *testing.T) {
	assert := assert.New(t)

	var didCall bool
	listener := NewHTTPRequestEventListener(func(_ context.Context, hre HTTPRequestEvent) {
		didCall = true
	})
	listener(context.Background(), logger.NewMessageEvent(logger.Info, "test"))
	assert.False(didCall)
	listener(context.Background(), NewHTTPRequestEvent(nil))
	assert.True(didCall)
}

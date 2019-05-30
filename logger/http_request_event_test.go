package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/blend/go-sdk/assert"
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
		OptHTTPRequestMeta(OptEventMetaFlag("test")),
		OptHTTPRequestRoute("/foo/:bar"),
		OptHTTPRequestState("this is the state"),
	)

	assert.NotNil(hre.Request)
	assert.Equal("GET", hre.Request.Method)
	assert.Equal("/foo", hre.Request.URL.Path)
	assert.Equal("test", hre.GetFlag())
	assert.Equal("/foo/:bar", hre.Route)
	assert.Equal("this is the state", hre.State)

	noColor := TextOutputFormatter{
		NoColor: true,
	}

	buf := new(bytes.Buffer)
	hre.WriteText(noColor, buf)
	assert.NotEmpty(buf.String())

	contents, err := json.Marshal(hre)
	assert.Nil(err)
	assert.NotEmpty(contents)
}

func TestHTTPRequestEventListener(t *testing.T) {
	assert := assert.New(t)

	var didCall bool
	listener := NewHTTPRequestEventListener(func(_ context.Context, hre *HTTPRequestEvent) {
		didCall = true
	})
	listener(context.Background(), NewMessageEvent(Info, "test"))
	assert.False(didCall)
	listener(context.Background(), NewHTTPRequestEvent(nil))
	assert.True(didCall)
}

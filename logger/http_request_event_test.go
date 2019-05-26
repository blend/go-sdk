package logger

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNewHTTPRequestEvent(t *testing.T) {
	assert := assert.New(t)

	hre := NewHTTPRequestEvent(nil,
		OptHTTPRequestEventRequest(&http.Request{
			Method: "GET",
			URL: &url.URL{
				Scheme: "https",
				Host:   "localhost",
				Path:   "/foo",
			},
		}),
		OptHTTPRequestEventOptionMetaOptions(OptEventMetaFlag("test")),
	)

	assert.NotNil(hre.Request)
	assert.Equal("test", hre.GetFlag())

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

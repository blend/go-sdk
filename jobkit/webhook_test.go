package jobkit

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestWebhook(t *testing.T) {
	assert := assert.New(t)

	wh := Webhook{
		Method: "POST",
		URL:    "https://example.org/foo?fuzz=buzz",
		Headers: map[string]string{
			"Authorization": "bailey",
		},
		Body: "this is a test",
	}

	request := wh.Request()
	assert.Equal("POST", request.Method)
	assert.NotNil(request.URL)
	assert.Equal("https", request.URL.Scheme)
	assert.Equal("example.org", request.URL.Host)
	assert.Equal("/foo", request.URL.Path)
	assert.NotEmpty(request.Header)
	assert.Equal("bailey", request.Header.Get("Authorization"))
	assert.NotNil(request.Body)
}

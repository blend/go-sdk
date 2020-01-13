package secrets

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestMockHTTPClientDo(t *testing.T) {
	assert := assert.New(t)

	// mocked request and response
	address := "www.blend.com"
	url, err := url.Parse(address)
	assert.Nil(err)

	client := NewMockHTTPClient()

	// Test: Do returns the OK response matching the route specified
	happyReq, err := http.NewRequest(http.MethodGet, address, bytes.NewReader([]byte{}))
	happyResp := &http.Response{StatusCode: http.StatusOK}
	assert.Nil(err)
	r, err := client.With(http.MethodGet, url, happyResp).Do(happyReq)
	assert.Nil(err)
	assert.Equal(happyResp, r)

	// Test: Do returns an error response for unknown routes
	badReq, err := http.NewRequest(http.MethodDelete, address, bytes.NewReader([]byte{}))
	r, err = client.Do(badReq)
	assert.NotNil(err)
	assert.Nil(r)
}

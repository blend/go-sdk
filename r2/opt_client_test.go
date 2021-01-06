package r2

import (
	"net/http"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestOptClient(t *testing.T) {
	it := assert.New(t)

	client := &http.Client{
		Timeout: time.Second,
	}
	req := New("https://localhost:8080", OptClient(client))
	it.Nil(req.Err)
	it.NotNil(req.Client)
	it.Equal(time.Second, req.Client.Timeout)

	xport := &http.Transport{}
	err := OptTransport(xport)(req)
	it.Nil(err)
	it.Equal(time.Second, req.Client.Timeout)
	it.NotNil(req.Client.Transport)
}

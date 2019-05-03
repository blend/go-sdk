package logger

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/webutil"
)

func TestHTTPLogged(t *testing.T) {
	assert := assert.New(t)

	buf := new(bytes.Buffer)
	log, err := New(
		OptOutput(buf),
		OptAll(),
	)
	assert.Nil(err)

	var didCall bool
	server := httptest.NewServer(webutil.NestMiddleware(func(rw http.ResponseWriter, req *http.Request) {
		didCall = true
	}, HTTPLogged(log)))

	res, err := http.Get(server.URL)
	assert.Nil(err)
	defer res.Body.Close()
	assert.True(didCall)
}

package webutil

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
)

func TestHTTPLogged(t *testing.T) {
	assert := assert.New(t)

	buf := new(bytes.Buffer)
	log, err := logger.New(
		logger.OptOutput(buf),
		logger.OptAll(),
	)
	assert.Nil(err)

	var didCall bool
	server := httptest.NewServer(NestMiddleware(func(rw http.ResponseWriter, req *http.Request) {
		didCall = true
	}, HTTPLogged(log)))

	res, err := http.Get(server.URL)
	assert.Nil(err)
	defer res.Body.Close()
	assert.True(didCall)
}

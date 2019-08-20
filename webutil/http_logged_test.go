package webutil

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blend/go-sdk/ansi"
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

func TestFormatHeaders(t *testing.T) {
	assert := assert.New(t)

	tf := logger.NewTextOutputFormatter(logger.OptTextNoColor())
	actual := FormatHeaders(tf, ansi.ColorBlue, http.Header{"Foo": []string{"bar"}, "Moo": []string{"loo"}})
	assert.Equal("{ Foo:bar Moo:loo }", actual)

	actual = FormatHeaders(tf, ansi.ColorBlue, http.Header{"Moo": []string{"loo"}, "Foo": []string{"bar"}})
	assert.Equal("{ Foo:bar Moo:loo }", actual)

	tf = logger.NewTextOutputFormatter()
	actual = FormatHeaders(tf, ansi.ColorBlue, http.Header{"Foo": []string{"bar"}, "Moo": []string{"loo"}})
	assert.Equal("{ "+ansi.ColorBlue.Apply("Foo")+":bar "+ansi.ColorBlue.Apply("Moo")+":loo }", actual)
}

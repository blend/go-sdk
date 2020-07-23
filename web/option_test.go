package web

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

func TestOptConfig(t *testing.T) {
	assert := assert.New(t)

	cfg := Config{
		DefaultHeaders: map[string]string{"X-Debug": "debug-value"},
		CookieName:     "FOOBAR",
	}

	var app App
	assert.Nil(OptConfig(cfg)(&app))
	assert.Equal("FOOBAR", app.Auth.CookieDefaults.Name)
	assert.NotEmpty(app.BaseHeaders)
	assert.Equal([]string{"debug-value"}, app.BaseHeaders["X-Debug"])
	assert.Equal([]string{PackageName}, app.BaseHeaders[webutil.HeaderServer])
}

func TestOptBindAddr(t *testing.T) {
	assert := assert.New(t)

	var app App
	assert.Nil(OptBindAddr(":9999")(&app))
	assert.Equal(":9999", app.Config.BindAddr)
}

func TestOptPort(t *testing.T) {
	assert := assert.New(t)

	var app App
	assert.Nil(OptPort(9999)(&app))
	assert.Equal(":9999", app.Config.BindAddr)
	assert.Equal(9999, app.Config.Port)
}

func TestOptLog(t *testing.T) {
	assert := assert.New(t)

	var app App
	assert.Nil(OptLog(logger.None())(&app))
	assert.NotNil(app.Log)
}

func TestOptServerOptions(t *testing.T) {
	t.Skip()

	assert := assert.New(t)

	baseline, baselineErr := New()
	assert.Nil(baselineErr)
	assert.NotNil(baseline)
	assert.NotNil(baseline.Server)
	assert.Nil(baseline.Server.ErrorLog)

	app, err := New(
		OptServerOptions(
			webutil.OptHTTPServerErrorLog(log.New(ioutil.Discard, "", log.LstdFlags)),
		),
	)
	assert.Nil(err)
	assert.NotNil(app.Server.ErrorLog)

	go func() { _ = app.Start() }()
	<-app.NotifyStarted()
	defer func() { _ = app.Stop() }()

	assert.NotNil(app.Server.ErrorLog)
}

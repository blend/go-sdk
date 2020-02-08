package web

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
)

func TestOptConfig(t *testing.T) {
	assert := assert.New(t)

	var app App
	assert.Nil(OptConfig(Config{CookieName: "FOOBAR"})(&app))
	assert.Equal("FOOBAR", app.Auth.CookieDefaults.Name)
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
	assert := assert.New(t)

	app, err := New(OptServerOptions(
		func(s *http.Server) error {
			s.ErrorLog = log.New(ioutil.Discard, "", log.LstdFlags)
			return nil
		},
	))
	assert.Nil(err)

	app.Start()
	defer app.Stop()
	<-app.NotifyStarted()
	assert.NotNil(app.Server.ErrorLog)
}

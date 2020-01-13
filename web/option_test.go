package web

import (
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

package web

import (
	"testing"

	"github.com/blend/go-sdk/logger"

	"github.com/blend/go-sdk/assert"
)

func TestOptConfig(t *testing.T) {
	assert := assert.New(t)

	var app App
	assert.Nil(OptConfig(Config{})(&app))
	assert.NotNil(app.Auth.FetchHandler)
}

func TestOptConfigError(t *testing.T) {
	assert := assert.New(t)

	var app App
	assert.NotNil(OptConfig(Config{AuthManagerMode: "NOT A REAL MODE"})(&app))
	assert.Nil(app.Auth.FetchHandler)
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

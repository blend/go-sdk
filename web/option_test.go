package web

import (
	"testing"

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

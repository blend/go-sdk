/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

import (
	"context"
	"io"
	"log"
	"net"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

func TestOptConfig(t *testing.T) {
	assert := assert.New(t)

	cfg := Config{
		DefaultHeaders:	map[string]string{"X-Debug": "debug-value"},
		CookieName:	"FOOBAR",
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

func TestOptBaseContext(t *testing.T) {
	t.Parallel()
	it := assert.New(t)

	ctx := context.TODO()
	opt := OptBaseContext(ctx)
	a := App{}
	err := opt(&a)
	it.Nil(err)
	root := a.BaseContext(nil)
	it.ReferenceEqual(ctx, root)
}

func TestOptBaseContextFunc(t *testing.T) {
	t.Parallel()
	it := assert.New(t)

	ctx := context.TODO()
	bc := func(_ net.Listener) context.Context {
		return ctx
	}
	opt := OptBaseContextFunc(bc)
	a := App{}
	err := opt(&a)
	it.Nil(err)
	root := a.BaseContext(nil)
	it.ReferenceEqual(ctx, root)
}

func TestOptLog(t *testing.T) {
	assert := assert.New(t)

	var app App
	assert.Nil(OptLog(logger.None())(&app))
	assert.NotNil(app.Log)
}

func TestOptServerOptions(t *testing.T) {
	assert := assert.New(t)

	baseline, baselineErr := New()
	assert.Nil(baselineErr)
	assert.NotNil(baseline)
	assert.NotNil(baseline.Server)
	assert.Nil(baseline.Server.ErrorLog)

	app, err := New(
		OptBindAddr("127.0.0.1:0"),
		OptServerOptions(
			webutil.OptHTTPServerErrorLog(log.New(io.Discard, "", log.LstdFlags)),
		),
	)
	assert.Nil(err)
	assert.NotNil(app.Server.ErrorLog)

	go func() { _ = app.Start() }()
	<-app.NotifyStarted()
	defer func() { _ = app.Stop() }()

	assert.NotNil(app.Server.ErrorLog)
}

func TestOptReadTimeout(t *testing.T) {
	assert := assert.New(t)

	var app App
	assert.Zero(app.Config.ReadTimeout)
	assert.Nil(OptReadTimeout(time.Second)(&app))
	assert.Equal(time.Second, app.Config.ReadTimeout)
}

func TestOptReadHeaderTimeout(t *testing.T) {
	assert := assert.New(t)

	var app App
	assert.Zero(app.Config.ReadHeaderTimeout)
	assert.Nil(OptReadHeaderTimeout(time.Second)(&app))
	assert.Equal(time.Second, app.Config.ReadHeaderTimeout)
}

func TestOptWriteTimeout(t *testing.T) {
	assert := assert.New(t)

	var app App
	assert.Zero(app.Config.WriteTimeout)
	assert.Nil(OptWriteTimeout(time.Second)(&app))
	assert.Equal(time.Second, app.Config.WriteTimeout)
}

func TestOptIdleTimeout(t *testing.T) {
	assert := assert.New(t)

	var app App
	assert.Zero(app.Config.IdleTimeout)
	assert.Nil(OptIdleTimeout(time.Second)(&app))
	assert.Equal(time.Second, app.Config.IdleTimeout)
}

func TestOptMaxHeaderBytes(t *testing.T) {
	assert := assert.New(t)

	var app App
	assert.Zero(app.Config.MaxHeaderBytes)
	assert.Nil(OptMaxHeaderBytes(100)(&app))
	assert.Equal(100, app.Config.MaxHeaderBytes)
}

func TestOptBaseURL(t *testing.T) {
	assert := assert.New(t)

	var app App
	assert.Empty(app.Config.BaseURL)
	assert.Nil(OptBaseURL("https://example.local")(&app))
	assert.Equal("https://example.local", app.Config.BaseURL)
}

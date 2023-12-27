/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"context"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/webutil"
)

var (
	_ configutil.Resolver = (*Config)(nil)
)

func TestConfigBindAddrOrDefault(t *testing.T) {
	assert := assert.New(t)
	var c Config
	assert.Equal(DefaultBindAddr, c.BindAddrOrDefault())
	assert.Equal("localhost:10", c.BindAddrOrDefault("localhost:10"))
	c.Port = 10
	assert.Equal(":10", c.BindAddrOrDefault())
	c.BindAddr = "localhost:10"
	assert.Equal(c.BindAddr, c.BindAddrOrDefault())
}

func TestConfigPortOrDefault(t *testing.T) {
	assert := assert.New(t)
	var c Config
	assert.Equal(webutil.PortFromBindAddr(DefaultBindAddr), c.PortOrDefault())
	c.BindAddr = ":10"
	assert.Equal(10, c.PortOrDefault())
	c.Port = 10
	assert.Equal(c.Port, c.PortOrDefault())
}

func TestConfigSessionTimeoutOrDefault(t *testing.T) {
	assert := assert.New(t)
	var c Config
	assert.Equal(DefaultSessionTimeout, c.SessionTimeoutOrDefault())
	c.SessionTimeout = 10
	assert.Equal(c.SessionTimeout, c.SessionTimeoutOrDefault())
}

func TestConfigCookieNameOrDefault(t *testing.T) {
	assert := assert.New(t)
	var c Config
	assert.Equal(DefaultCookieName, c.CookieNameOrDefault())
	c.CookieName = "helloworld"
	assert.Equal(c.CookieName, c.CookieNameOrDefault())
}

func TestConfigCookiePathOrDefault(t *testing.T) {
	assert := assert.New(t)
	var c Config
	assert.Equal(DefaultCookiePath, c.CookiePathOrDefault())
	c.CookiePath = "helloworld"
	assert.Equal(c.CookiePath, c.CookiePathOrDefault())
}

func TestConfigCookieSecureOrDefault(t *testing.T) {
	assert := assert.New(t)
	var c Config
	yes := true
	assert.Equal(DefaultCookieSecure, c.CookieSecureOrDefault())
	c.BaseURL = "https://hello.com"
	assert.True(c.CookieSecureOrDefault())
	c.BaseURL = "http://hello.com"
	assert.False(c.CookieSecureOrDefault())
	c.CookieSecure = &yes
	assert.Equal(*c.CookieSecure, c.CookieSecureOrDefault())
}

func TestConfigCookieHTTPOnlyOrDefault(t *testing.T) {
	assert := assert.New(t)
	var c Config
	yes := true
	assert.Equal(DefaultCookieHTTPOnly, c.CookieHTTPOnlyOrDefault())
	c.CookieHTTPOnly = &yes
	assert.Equal(*c.CookieHTTPOnly, c.CookieHTTPOnlyOrDefault())
}

func TestConfigCookieSameSiteOrDefault(t *testing.T) {
	assert := assert.New(t)
	var c Config
	assert.Equal(DefaultCookieSameSiteMode, c.CookieSameSiteOrDefault())

	c.CookieSameSite = webutil.SameSiteStrict
	assert.Equal(http.SameSiteStrictMode, c.CookieSameSiteOrDefault())

	assert.NotNil(func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = ex.New(r)
			}
		}()
		c.CookieSameSite = "not valid"
		assert.Equal(c.CookieSameSite, c.CookieSameSiteOrDefault())
		return
	}())
}

func TestConfigMaxHeaderBytesOrDefault(t *testing.T) {
	assert := assert.New(t)
	var c Config
	assert.Equal(DefaultMaxHeaderBytes, c.MaxHeaderBytesOrDefault())
	c.MaxHeaderBytes = 1000
	assert.Equal(c.MaxHeaderBytes, c.MaxHeaderBytesOrDefault())
}

func TestConfigReadTimeoutOrDefault(t *testing.T) {
	assert := assert.New(t)
	var c Config
	assert.Equal(DefaultReadTimeout, c.ReadTimeoutOrDefault())
	c.ReadTimeout = 1000
	assert.Equal(c.ReadTimeout, c.ReadTimeoutOrDefault())
}

func TestConfigReadHeaderTimeoutOrDefault(t *testing.T) {
	assert := assert.New(t)
	var c Config
	assert.Equal(DefaultReadHeaderTimeout, c.ReadHeaderTimeoutOrDefault())
	c.ReadHeaderTimeout = 1000
	assert.Equal(c.ReadHeaderTimeout, c.ReadHeaderTimeoutOrDefault())
}

func TestConfigWriteTimeoutOrDefault(t *testing.T) {
	assert := assert.New(t)
	var c Config
	assert.Equal(DefaultWriteTimeout, c.WriteTimeoutOrDefault())
	c.WriteTimeout = 1000
	assert.Equal(c.WriteTimeout, c.WriteTimeoutOrDefault())
}

func TestConfigIdleTimeoutOrDefault(t *testing.T) {
	assert := assert.New(t)
	var c Config
	assert.Equal(DefaultIdleTimeout, c.IdleTimeoutOrDefault())
	c.IdleTimeout = 1000
	assert.Equal(c.IdleTimeout, c.IdleTimeoutOrDefault())
}

func TestConfigShutdownGracePeriodOrDefault(t *testing.T) {
	assert := assert.New(t)
	var c Config
	assert.Equal(DefaultShutdownGracePeriod, c.ShutdownGracePeriodOrDefault())
	c.ShutdownGracePeriod = 1000
	assert.Equal(c.ShutdownGracePeriod, c.ShutdownGracePeriodOrDefault())
}

func TestConfigResolve(t *testing.T) {
	assert := assert.New(t)

	var c Config
	env.SetEnv(env.New())

	defer env.Restore()
	assert.Nil(c.Resolve(env.WithVars(context.Background(), env.Env())))
	assert.Empty(c.BindAddr)

	env.Env().Set("BIND_ADDR", "hello")
	assert.Nil(c.Resolve(env.WithVars(context.Background(), env.Env())))
	assert.Equal("hello", c.BindAddr)
}

func TestConfigResolve_CookieDomainFromBaseURL(t *testing.T) {
	its := assert.New(t)

	cfg := Config{
		BaseURL: "https://example.com/foo/bar?buzz=fuzz",
	}
	err := (&cfg).Resolve(context.TODO())
	its.Nil(err)
	its.True(*cfg.CookieSecure)
	its.Equal("example.com", cfg.CookieDomain)
}

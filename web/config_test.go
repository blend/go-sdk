package web

import (
	"encoding/base64"
	"testing"

	"github.com/blend/go-sdk/env"

	"github.com/blend/go-sdk/webutil"

	"github.com/blend/go-sdk/assert"
)

func TestConfigBindAddrOrDefault(t *testing.T) {
	assert := assert.New(t)
	c := &Config{}
	assert.Equal(DefaultBindAddr, c.BindAddrOrDefault())
	assert.Equal("localhost:10", c.BindAddrOrDefault("localhost:10"))
	c.Port = 10
	assert.Equal(":10", c.BindAddrOrDefault())
	c.BindAddr = "localhost:10"
	assert.Equal(c.BindAddr, c.BindAddrOrDefault())
}

func TestConfigPortOrDefault(t *testing.T) {
	assert := assert.New(t)
	c := &Config{}
	assert.Equal(webutil.PortFromBindAddr(DefaultBindAddr), c.PortOrDefault())
	c.BindAddr = ":10"
	assert.Equal(10, c.PortOrDefault())
	c.Port = 10
	assert.Equal(c.Port, c.PortOrDefault())
}

func TestConfigBaseURLIsSecureScheme(t *testing.T) {
	assert := assert.New(t)
	c := &Config{}
	assert.False(c.BaseURLIsSecureScheme())
	c.BaseURL = "http://hello.com"
	assert.False(c.BaseURLIsSecureScheme())
	c.BaseURL = "https://hello.com"
	assert.True(c.BaseURLIsSecureScheme())
}

func TestConfigAuthManagerModeOrDefault(t *testing.T) {
	assert := assert.New(t)
	c := &Config{}
	assert.Equal(AuthManagerModeRemote, c.AuthManagerModeOrDefault())
	c.AuthManagerMode = string(AuthManagerModeJWT)
	assert.Equal(c.AuthManagerMode, c.AuthManagerModeOrDefault())
}

func TestConfigSessionTimeoutOrDefault(t *testing.T) {
	assert := assert.New(t)
	c := &Config{}
	assert.Equal(DefaultSessionTimeout, c.SessionTimeoutOrDefault())
	c.SessionTimeout = 10
	assert.Equal(c.SessionTimeout, c.SessionTimeoutOrDefault())
}

func TestConfigCookieNameOrDefault(t *testing.T) {
	assert := assert.New(t)
	c := &Config{}
	assert.Equal(DefaultCookieName, c.CookieNameOrDefault())
	c.CookieName = "helloworld"
	assert.Equal(c.CookieName, c.CookieNameOrDefault())
}

func TestConfigCookiePathOrDefault(t *testing.T) {
	assert := assert.New(t)
	c := &Config{}
	assert.Equal(DefaultCookiePath, c.CookiePathOrDefault())
	c.CookiePath = "helloworld"
	assert.Equal(c.CookiePath, c.CookiePathOrDefault())
}

func TestConfigCookieSecureOrDefault(t *testing.T) {
	assert := assert.New(t)
	c := &Config{}
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
	c := &Config{}
	yes := true
	assert.Equal(DefaultCookieHTTPOnly, c.CookieHTTPOnlyOrDefault())
	c.CookieHTTPOnly = &yes
	assert.Equal(*c.CookieHTTPOnly, c.CookieHTTPOnlyOrDefault())
}

func TestConfigCookieSameSiteOrDefault(t *testing.T) {
	assert := assert.New(t)
	c := &Config{}
	assert.Equal(DefaultCookieSameSite, c.CookieSameSiteOrDefault())
	c.CookieSameSite = "helloworld"
	assert.Equal(c.CookieSameSite, c.CookieSameSiteOrDefault())
}

func TestConfigMaxHeaderBytesOrDefault(t *testing.T) {
	assert := assert.New(t)
	c := &Config{}
	assert.Equal(DefaultMaxHeaderBytes, c.MaxHeaderBytesOrDefault())
	c.MaxHeaderBytes = 1000
	assert.Equal(c.MaxHeaderBytes, c.MaxHeaderBytesOrDefault())
}

func TestConfigReadTimeoutOrDefault(t *testing.T) {
	assert := assert.New(t)
	c := &Config{}
	assert.Equal(DefaultReadTimeout, c.ReadTimeoutOrDefault())
	c.ReadTimeout = 1000
	assert.Equal(c.ReadTimeout, c.ReadTimeoutOrDefault())
}

func TestConfigReadHeaderTimeoutOrDefault(t *testing.T) {
	assert := assert.New(t)
	c := &Config{}
	assert.Equal(DefaultReadHeaderTimeout, c.ReadHeaderTimeoutOrDefault())
	c.ReadHeaderTimeout = 1000
	assert.Equal(c.ReadHeaderTimeout, c.ReadHeaderTimeoutOrDefault())
}

func TestConfigWriteTimeoutOrDefault(t *testing.T) {
	assert := assert.New(t)
	c := &Config{}
	assert.Equal(DefaultWriteTimeout, c.WriteTimeoutOrDefault())
	c.WriteTimeout = 1000
	assert.Equal(c.WriteTimeout, c.WriteTimeoutOrDefault())
}

func TestConfigIdleTimeoutOrDefault(t *testing.T) {
	assert := assert.New(t)
	c := &Config{}
	assert.Equal(DefaultIdleTimeout, c.IdleTimeoutOrDefault())
	c.IdleTimeout = 1000
	assert.Equal(c.IdleTimeout, c.IdleTimeoutOrDefault())
}

func TestConfigShutdownGracePeriodOrDefault(t *testing.T) {
	assert := assert.New(t)
	c := &Config{}
	assert.Equal(DefaultShutdownGracePeriod, c.ShutdownGracePeriodOrDefault())
	c.ShutdownGracePeriod = 1000
	assert.Equal(c.ShutdownGracePeriod, c.ShutdownGracePeriodOrDefault())
}

func TestConfigMustAuthSecret(t *testing.T) {
	assert := assert.New(t)
	c := &Config{}
	s := []byte("secret")
	c.AuthSecret = base64.StdEncoding.EncodeToString(s)
	assert.Equal(string(s), string(c.MustAuthSecret()))

	c.AuthSecret = "non decodable"
	panicy := func() {
		defer func() {
			err := recover()
			assert.NotNil(err)
		}()
		c.MustAuthSecret()
	}
	panicy()
}

func TestConfigResolve(t *testing.T) {
	assert := assert.New(t)
	c := &Config{}
	env.SetEnv(env.New())
	assert.Nil(c.Resolve())
	assert.Empty(c.BindAddr)
	env.Env().Set("BIND_ADDR", "hello")
	assert.Nil(c.Resolve())
	assert.Equal("hello", c.BindAddr)
}

package web

import (
	"testing"

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

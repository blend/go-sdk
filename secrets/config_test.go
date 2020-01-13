package secrets

import (
	"testing"
	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
)

func TestNewConfigFromEnv(t *testing.T) {
	assert := assert.New(t)
	defer env.Restore()

	env.Env().Set("VAULT_ADDR", "http://127.0.0.2:8100")
	env.Env().Set("VAULT_TOKEN", "thisisatest")

	cfg, err := NewConfigFromEnv()
	assert.Nil(err)
	assert.Equal("http://127.0.0.2:8100", cfg.AddrOrDefault())
	assert.Equal("thisisatest", cfg.Token)
}

func TestConfigIsZero(t *testing.T) {
	assert := assert.New(t)

	assert.True(Config{}.IsZero())
	assert.False(Config{Token: "garbage"}.IsZero())
}

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	cfg := Config{}
	assert.Equal(DefaultAddr, cfg.AddrOrDefault())
	assert.Empty(cfg.Token)
	assert.Equal(DefaultMount, cfg.MountOrDefault())
	assert.Equal(DefaultTimeout, cfg.TimeoutOrDefault())
	assert.Empty(cfg.RootCAs)
	assert.Empty(cfg.ServicePath)
}

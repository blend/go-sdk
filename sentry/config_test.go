package sentry

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
)

func TestConfigResolve(t *testing.T) {
	assert := assert.New(t)

	env.Restore()
	env.Env().Set("SENTRY_DSN", "http://foo@example.com/1")
	env.Env().Set(env.VarServiceName, "go-sdk-server")
	env.Env().Set("SENTRY_SERVER_NAME", "go-sdk-web-server")
	env.Env().Set("SENTRY_DIST", "v1.0.0")
	env.Env().Set("SENTRY_RELEASE", "deadbeef")
	env.Env().Set(env.VarServiceEnv, "dev")
	env.Env().Set("SENTRY_ENVIRONMENT", "test")

	cfg := &Config{}
	assert.True(cfg.IsZero())
	assert.Equal("go-sdk-server", cfg.ServerNameOrDefault())
	assert.Equal("dev", cfg.EnvironmentOrDefault())
	assert.Empty(cfg.DistOrDefault())
	assert.Empty(cfg.ReleaseOrDefault())

	assert.Nil(cfg.Resolve())
	assert.False(cfg.IsZero())
	assert.Equal("http://foo@example.com/1", cfg.DSN)
	assert.Equal("go-sdk-web-server", cfg.ServerName)
	assert.Equal("v1.0.0", cfg.Dist)
	assert.Equal("deadbeef", cfg.Release)
	assert.Equal("test", cfg.Environment)
}

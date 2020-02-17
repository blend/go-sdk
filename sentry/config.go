package sentry

import (
	"context"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/env"
)

// Config is the sentry config.
type Config struct {
	// The DSN to use. If the DSN is not set, the client is effectively disabled.
	DSN string `json:"dsn" yaml:"dsn" env:"SENTRY_DSN"`
	// The server name to be reported.
	ServerName string `json:"serverName" yaml:"serverName" env:"SENTRY_SERVER_NAME"`
	// The dist to be sent with events.
	Dist string `json:"dist" yaml:"dist" env:"SENTRY_DIST"`
	// The release to be sent with events.
	Release string `json:"release" yaml:"release" env:"SENTRY_RELEASE"`
	// The environment to be sent with events.
	Environment string `json:"environment" yaml:"environment" env:"SENTRY_ENVIRONMENT"`
	// Maximum number of breadcrumbs.
	MaxBreadcrumbs int `json:"maxBreadCrumbs" yaml:"maxBreadCrumbs"`
}

// IsZero returns if the config is unset.
func (c Config) IsZero() bool {
	return c.DSN == ""
}

// Resolve applies configutil resoltion steps.
func (c *Config) Resolve(ctx context.Context) error {
	return configutil.GetEnvVars(ctx).ReadInto(c)
}

// ServerNameOrDefault returns the server name or a default.
func (c Config) ServerNameOrDefault() string {
	if c.ServerName != "" {
		return c.ServerName
	}
	return env.Env().ServiceName()
}

// EnvironmentOrDefault returns the environment or a default.
func (c Config) EnvironmentOrDefault() string {
	if c.Environment != "" {
		return c.Environment
	}
	return env.Env().ServiceEnv()
}

// DistOrDefault returns the dist or a default.
func (c Config) DistOrDefault() string {
	if c.Dist != "" {
		return c.Dist
	}
	return ""
}

// ReleaseOrDefault returns the release or a default.
func (c Config) ReleaseOrDefault() string {
	if c.Release != "" {
		return c.Release
	}
	return ""
}

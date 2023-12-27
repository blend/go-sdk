/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sentry

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/env"
)

// Config is the sentry config.
type Config struct {
	// The DSN to use. If the DSN is not set, the client is effectively disabled.
	DSN string `json:"dsn" yaml:"dsn"`
	// The server name to be reported.
	ServerName string `json:"serverName" yaml:"serverName"`
	// The dist to be sent with events.
	Dist string `json:"dist" yaml:"dist"`
	// The release to be sent with events.
	Release string `json:"release" yaml:"release"`
	// The environment to be sent with events.
	Environment string `json:"environment" yaml:"environment"`
	// Maximum number of breadcrumbs.
	MaxBreadcrumbs int `json:"maxBreadCrumbs" yaml:"maxBreadCrumbs"`
	// Debug prints debugging information to the screen.
	Debug bool `json:"debug" yaml:"debug"`
	// FlushTimeout is the timeout for flushing exceptions to sentry.
	FlushTimeout time.Duration `json:"flushTimeout" yaml:"flushTimeout"`
}

// IsZero returns if the config is unset.
func (c Config) IsZero() bool {
	return c.DSN == ""
}

// Resolve applies configutil resoltion steps.
func (c *Config) Resolve(ctx context.Context) error {
	return configutil.Resolve(ctx,
		configutil.SetString(&c.DSN, configutil.String(c.DSN), configutil.Env("SENTRY_DSN")),
		configutil.SetString(&c.ServerName, configutil.String(c.ServerName), configutil.Env(env.VarServiceName)),
		configutil.SetString(&c.Environment, configutil.String(c.Environment), configutil.Env(env.VarServiceEnv)),
		configutil.SetDuration(&c.FlushTimeout, configutil.Duration(c.FlushTimeout), configutil.Duration(5*time.Second)),
	)
}

// GetDSNHost returns just the scheme and hostname for the dsn.
func (c *Config) GetDSNHost() string {
	if c.DSN == "" {
		return ""
	}

	parsedURL, _ := url.Parse(c.DSN)
	if parsedURL == nil {
		return ""
	}
	return fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
}

// FlushTimeoutOrDefault returns the flush timeout or a default.
func (c Config) FlushTimeoutOrDefault() time.Duration {
	if c.FlushTimeout > 0 {
		return c.FlushTimeout
	}
	return 5 * time.Second
}

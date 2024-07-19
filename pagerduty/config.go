/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package pagerduty

import (
	"context"

	"github.com/blend/go-sdk/configutil"
)

// Config is the pagerduty config.
type Config struct {
	Addr  string `yaml:"addr,omitempty"`
	Token string `yaml:"token,omitempty"`
	Email string `yaml:"email,omitempty"`
}

// IsZero returns if the config is set or not.
func (c Config) IsZero() bool {
	return c.Token == "" || c.Email == ""
}

// Resolve resolves the config.
func (c *Config) Resolve(ctx context.Context) error {
	return configutil.Resolve(ctx,
		configutil.SetString(&c.Addr, configutil.String(c.Addr), configutil.String(DefaultAddr)),
	)
}

/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package vault

import (
	"context"
	"time"

	"github.com/blend/go-sdk/configutil"
)

// Config is the secrets config object.
type Config struct {
	// Addr is the remote address of the secret store.
	Addr	string	`json:"addr" yaml:"addr" env:"VAULT_ADDR"`
	// Mount is the default mount path, it prefixes any paths.
	Mount	string	`json:"mount" yaml:"mount" env:"VAULT_MOUNT"`
	// Token is the authentication token used to talk to the secret store.
	Token	string	`json:"token" yaml:"token" env:"VAULT_TOKEN"`
	// Timeout is the dial timeout for requests to the secrets store.
	Timeout	time.Duration	`json:"timeout" yaml:"timeout" env:"VAULT_TIMEOUT"`
	// RootCAs is a list of certificate authority paths.
	RootCAs	[]string	`json:"rootCAs" yaml:"rootCAs" env:"VAULT_CA_CERT,csv"`
}

// IsZero returns if the config is set or not.
func (c Config) IsZero() bool {
	return len(c.Token) == 0
}

// Resolve reads the environment into the config on configutil.Read(...)
func (c *Config) Resolve(ctx context.Context) error {
	return configutil.Resolve(ctx,
		configutil.SetString(&c.Addr, configutil.String(c.Addr), configutil.Env(EnvVarVaultAddr)),
		configutil.SetString(&c.Mount, configutil.String(c.Mount), configutil.Env(EnvVarVaultMount)),
		configutil.SetString(&c.Token, configutil.String(c.Token), configutil.Env(EnvVarVaultToken)),
		configutil.SetStrings(&c.RootCAs, configutil.Strings(c.RootCAs), configutil.Env(EnvVarVaultCertAuthorityPath)),
		configutil.SetDuration(&c.Timeout, configutil.Duration(c.Timeout), configutil.Env(EnvVarVaultTimeout)),
	)
}

// AddrOrDefault returns the client addr.
func (c Config) AddrOrDefault() string {
	if c.Addr != "" {
		return c.Addr
	}
	return DefaultAddr
}

// TimeoutOrDefault returns the client timeout.
func (c Config) TimeoutOrDefault() time.Duration {
	if c.Timeout > 0 {
		return c.Timeout
	}
	return DefaultTimeout
}

// MountOrDefault returns secrets mount or a default.
func (c Config) MountOrDefault() string {
	if c.Mount != "" {
		return c.Mount
	}
	return DefaultMount
}

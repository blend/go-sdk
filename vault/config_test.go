/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package vault

import (
	"context"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/env"
)

var (
	_ configutil.Resolver = (*Config)(nil)
)

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
	assert.Equal(DefaultTimeout, cfg.TimeoutOrDefault())
	assert.Empty(cfg.RootCAs)
}

func TestResolveTimeout(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	const fiveSeconds = "5s"
	defer env.Env().Restore(EnvVarVaultTimeout)

	env.Env().Set(EnvVarVaultTimeout, fiveSeconds)
	cfg := &Config{}
	err := cfg.Resolve(ctx)
	assert.Nil(err)
	assert.Equal(time.Second*5, cfg.TimeoutOrDefault())

	env.Env().Delete(EnvVarVaultTimeout)
	cfg2 := &Config{}
	err = cfg2.Resolve(ctx)
	assert.Equal(DefaultTimeout, cfg2.TimeoutOrDefault())
	assert.Nil(err)
}

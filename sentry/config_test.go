/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sentry

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
)

func Test_Config_Resolve(t *testing.T) {
	its := assert.New(t)

	vars := env.Vars{
		"SENTRY_DSN":		"test-dsn",
		env.VarServiceEnv:	env.ServiceEnvTest,
		env.VarServiceName:	"sentry-test",
	}

	cfg := new(Config)
	err := cfg.Resolve(env.WithVars(context.Background(), vars))
	its.Nil(err)
	its.Equal("test-dsn", cfg.DSN)
	its.False(cfg.IsZero())
	its.Equal(env.ServiceEnvTest, cfg.Environment)
	its.Equal("sentry-test", cfg.ServerName)
}

func Test_Config_Resolve_noDSN(t *testing.T) {
	its := assert.New(t)

	vars := env.Vars{
		env.VarServiceEnv:	env.ServiceEnvTest,
		env.VarServiceName:	"sentry-test",
	}

	cfg := new(Config)
	err := cfg.Resolve(env.WithVars(context.Background(), vars))
	its.Nil(err)
	its.Empty(cfg.DSN)
	its.True(cfg.IsZero())
	its.Equal(env.ServiceEnvTest, cfg.Environment)
	its.Equal("sentry-test", cfg.ServerName)
}

func Test_Config_GetDSNHost(t *testing.T) {
	its := assert.New(t)

	cfg := &Config{
		DSN: "https://admin:nopasswd@example.com/buzz/fuzz?query=value",
	}
	its.Equal("https://example.com", cfg.GetDSNHost())
}

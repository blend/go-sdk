/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
)

func TestEventFlagSetEnable(t *testing.T) {
	assert := assert.New(t)

	set := NewFlags("FOO")
	set.Enable("TEST")
	assert.False(set.IsEnabled("TEST"))
	assert.False(set.IsEnabled("FOO"))
	assert.True(set.IsEnabled("test"))
	assert.True(set.IsEnabled("foo"))
	assert.False(set.IsEnabled("NOT_TEST"))
}

func TestEventFlagSetDisable(t *testing.T) {
	assert := assert.New(t)

	set := NewFlags()
	set.Enable("TEST")
	assert.True(set.IsEnabled("test"))
	set.Disable("TEST")
	assert.False(set.IsEnabled("test"))
}

func TestEventFlagSetEnableAll(t *testing.T) {
	assert := assert.New(t)

	set := NewFlags()
	set.SetAll()
	assert.True(set.IsEnabled("test"))
	assert.True(set.IsEnabled("NOT_TEST"))
	assert.True(set.IsEnabled("NOT_TEST"))
	set.Disable("test")
	assert.True(set.IsEnabled("NOT_TEST"))
	assert.False(set.IsEnabled("test"))
}

func TestEventFlagSetFromEnvironment(t *testing.T) {
	assert := assert.New(t)

	env.Env().Set(EnvVarFlags, "error,info,web.request")
	defer env.Env().Restore(EnvVarFlags)

	set := NewFlags(env.Env().CSV(EnvVarFlags)...)
	assert.True(set.IsEnabled(Error))
	assert.True(set.IsEnabled(Info))
	assert.False(set.IsEnabled(Fatal))
}

func TestEventFlagSetFromEnvironmentWithDisabled(t *testing.T) {
	assert := assert.New(t)

	env.Env().Set(EnvVarFlags, "all,-debug")
	defer env.Env().Restore(EnvVarFlags)

	set := NewFlags(env.Env().CSV(EnvVarFlags)...)
	assert.True(set.IsEnabled(Error))
	assert.True(set.IsEnabled(Fatal))
	assert.True(set.IsEnabled("foo"))
	assert.False(set.IsEnabled(Debug))
}

func TestEventFlagSetFromEnvironmentAll(t *testing.T) {
	assert := assert.New(t)

	env.Env().Set(EnvVarFlags, "all")
	defer env.Env().Restore(EnvVarFlags)

	set := NewFlags(env.Env().CSV(EnvVarFlags)...)
	assert.True(set.All())
	assert.False(set.None())
	assert.True(set.IsEnabled(Error))
}

func TestEventFlagSetFromEnvNone(t *testing.T) {
	assert := assert.New(t)

	env.Env().Set(EnvVarFlags, "none")
	defer env.Env().Restore(EnvVarFlags)

	set := NewFlags(env.Env().CSV(EnvVarFlags)...)
	assert.False(set.All())
	assert.True(set.None())
	assert.False(set.IsEnabled(Error))
}

func TestEventFlagNoneEnableEvents(t *testing.T) {
	assert := assert.New(t)

	flags := FlagsNone()
	assert.False(flags.IsEnabled("test_flag"))

	flags.Enable("test_flag")
	assert.True(flags.IsEnabled("test_flag"))
}

func TestEventSetCoalesceWith(t *testing.T) {
	assert := assert.New(t)

	first := NewFlags(Info)
	first.MergeWith(NewFlags(Warning))
	assert.True(first.IsEnabled(Info))
	assert.True(first.IsEnabled(Warning))
	assert.False(first.IsEnabled(Fatal))

	second := NewFlags(Info)
	second.MergeWith(NewFlags("-info"))
	assert.False(second.IsEnabled(Info))
}

func TestFlagSetNone(t *testing.T) {
	assert := assert.New(t)
	assert.True(FlagsNone().None())
}

func TestFlagSetString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("none", FlagsNone().String())
	assert.Equal("all", FlagsAll().String())

	fs := NewFlags(Info, Debug, Error)
	assert.Contains(fs.String(), "info")
	assert.Contains(fs.String(), "debug")
	assert.Contains(fs.String(), "error")

	nfs := FlagsAll()
	nfs.Disable(Fatal)
	assert.Equal("all, -fatal", nfs.String())
}

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package profanity

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ref"
)

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	cfg := Config{}
	assert.False(cfg.DebugOrDefault())
	cfg.Debug = ref.Bool(true)
	assert.True(cfg.DebugOrDefault())

	assert.False(cfg.VerboseOrDefault())
	cfg.Verbose = ref.Bool(true)
	assert.True(cfg.VerboseOrDefault())

	assert.False(cfg.ExitFirstOrDefault())
	cfg.ExitFirst = ref.Bool(true)
	assert.True(cfg.ExitFirstOrDefault())

	assert.Equal(DefaultRulesFile, cfg.RulesFileOrDefault())
	cfg.RulesFile = "foo"
	assert.Equal("foo", cfg.RulesFileOrDefault())
}

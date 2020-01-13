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

	assert.False(cfg.FailFastOrDefault())
	cfg.FailFast = ref.Bool(true)
	assert.True(cfg.FailFastOrDefault())

	assert.Equal(DefaultRulesFile, cfg.RulesFileOrDefault())
	cfg.RulesFile = "foo"
	assert.Equal("foo", cfg.RulesFileOrDefault())
}

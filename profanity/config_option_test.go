package profanity

import (
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestConfigOptions(t *testing.T) {
	assert := assert.New(t)

	cfg := &Config{}

	assert.False(cfg.VerboseOrDefault())
	OptVerbose(true)(cfg)
	assert.True(cfg.VerboseOrDefault())

	assert.False(cfg.DebugOrDefault())
	OptDebug(true)(cfg)
	assert.True(cfg.DebugOrDefault())

	assert.False(cfg.FailFastOrDefault())
	OptFailFast(true)(cfg)
	assert.True(cfg.FailFastOrDefault())

	assert.Empty(cfg.Root)
	OptRoot("../foo")(cfg)
	assert.Equal("../foo", cfg.Root)

	assert.Equal(DefaultRulesFile, cfg.RulesFileOrDefault())
	OptRulesFile("my_rules.yml")(cfg)

	assert.Empty(cfg.Include)
	OptInclude("foo", "bar", "baz")(cfg)
	assert.Equal([]string{"foo", "bar", "baz"}, cfg.Include)

	assert.Empty(cfg.Exclude)
	OptExclude("foo", "bar", "baz")(cfg)
	assert.Equal([]string{"foo", "bar", "baz"}, cfg.Exclude)
}

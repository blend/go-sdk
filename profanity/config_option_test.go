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

	assert.Empty(cfg.Path)
	OptPath("../foo")(cfg)
	assert.Equal("../foo", cfg.Path)

	assert.Equal(DefaultRulesFile, cfg.RulesFileOrDefault())
	OptRulesFile("my_rules.yml")(cfg)

	assert.Empty(cfg.Files.Include)
	OptFilesInclude("foo", "bar", "baz")(cfg)
	assert.Equal([]string{"foo", "bar", "baz"}, cfg.Files.Include)

	assert.Empty(cfg.Files.Exclude)
	OptFilesExclude("foo", "bar", "baz")(cfg)
	assert.Equal([]string{"foo", "bar", "baz"}, cfg.Files.Exclude)
}

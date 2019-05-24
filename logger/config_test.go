package logger

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	var cfg Config
	assert.Equal(DefaultFlags, cfg.FlagsOrDefault())
	assert.Equal(FormatText, cfg.FormatOrDefault())
	_, ok := cfg.Formatter().(*TextOutputFormatter)
	assert.True(ok)

	cfg = Config{
		Flags:  []string{Info, Error},
		Format: FormatJSON,
	}

	assert.Equal([]string{Info, Error}, cfg.FlagsOrDefault())
	assert.Equal(FormatJSON, cfg.FormatOrDefault())
}

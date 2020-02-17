package logger

import (
	"context"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/env"
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

func TestConfigResolve(t *testing.T) {
	assert := assert.New(t)

	defer env.Restore()
	env.Env().Set("LOG_FLAGS", "foo,bar")
	env.Env().Set("LOG_HIDE_TIMESTAMP", "true")
	env.Env().Set("LOG_HIDE_FIELDS", "true")
	env.Env().Set("LOG_TIME_FORMAT", time.Kitchen)
	env.Env().Set("NO_COLOR", "true")

	cfg := &Config{}
	ctx := configutil.WithEnvVars(context.Background(), env.Env())
	assert.Nil(cfg.Resolve(ctx))

	assert.Any(cfg.Flags, func(v interface{}) bool { return v.(string) == "foo" })
	assert.Any(cfg.Flags, func(v interface{}) bool { return v.(string) == "bar" })
	assert.True(cfg.Text.HideTimestamp)
	assert.True(cfg.Text.HideFields)
	assert.True(cfg.Text.NoColor)
	assert.Equal(time.Kitchen, cfg.Text.TimeFormat)
}

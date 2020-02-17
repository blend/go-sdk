package logger

import (
	"context"
	"strings"

	"github.com/blend/go-sdk/configutil"
)

// Config is the logger config.
type Config struct {
	Flags  []string   `json:"flags,omitempty" yaml:"flags,omitempty" env:"LOG_FLAGS,csv"`
	Format string     `json:"format,omitempty" yaml:"format,omitempty" env:"LOG_FORMAT"`
	Text   TextConfig `json:"text,omitempty" yaml:"text,omitempty"`
	JSON   JSONConfig `json:"json,omitempty" yaml:"json,omitempty"`
}

// Resolve resolves the config.
func (c *Config) Resolve(ctx context.Context) error {
	return configutil.ReturnFirst(
		configutil.GetEnvVars(ctx).ReadInto(c),
	)
}

// FlagsOrDefault returns the enabled logger events.
func (c Config) FlagsOrDefault() []string {
	if len(c.Flags) > 0 {
		return c.Flags
	}
	return DefaultFlags
}

// FormatOrDefault returns the output format or a default.
func (c Config) FormatOrDefault() string {
	if c.Format != "" {
		return c.Format
	}
	return FormatText
}

// Formatter returns the configured writers
func (c Config) Formatter() WriteFormatter {
	switch strings.ToLower(string(c.FormatOrDefault())) {
	case FormatJSON:
		return NewJSONOutputFormatter(OptJSONConfig(c.JSON))
	case FormatText:
		return NewTextOutputFormatter(OptTextConfig(c.Text))
	default:
		return NewTextOutputFormatter(OptTextConfig(c.Text))
	}
}

// TextConfig is the config for a text formatter.
type TextConfig struct {
	HideTimestamp bool   `json:"hideTimestamp,omitempty" yaml:"hideTimestamp,omitempty" env:"LOG_HIDE_TIMESTAMP"`
	HideFields    bool   `json:"hideFields,omitempty" yaml:"hideFields,omitempty" env:"LOG_HIDE_FIELDS"`
	NoColor       bool   `json:"noColor,omitempty" yaml:"noColor,omitempty" env:"NO_COLOR"`
	TimeFormat    string `json:"timeFormat,omitempty" yaml:"timeFormat,omitempty" env:"LOG_TIME_FORMAT"`
}

// TimeFormatOrDefault returns a field value or a default.
func (twc TextConfig) TimeFormatOrDefault() string {
	if len(twc.TimeFormat) > 0 {
		return twc.TimeFormat
	}
	return DefaultTextTimeFormat
}

// JSONConfig is the config for a json formatter.
type JSONConfig struct {
	Pretty       bool   `json:"pretty,omitempty" yaml:"pretty,omitempty" env:"LOG_JSON_PRETTY"`
	PrettyPrefix string `json:"prettyPrefix,omitempty" yaml:"prettyPrefix,omitempty" env:"LOG_JSON_PRETTY_PREFIX"`
	PrettyIndent string `json:"prettyIndent,omitempty" yaml:"prettyIndent,omitempty" env:"LOG_JSON_PRETTY_INDENT"`
}

// PrettyPrefixOrDefault returns the pretty prefix or a default.
func (jc JSONConfig) PrettyPrefixOrDefault() string {
	return jc.PrettyPrefix
}

// PrettyIndentOrDefault returns the pretty indent or a default.
func (jc JSONConfig) PrettyIndentOrDefault() string {
	if jc.PrettyIndent != "" {
		return jc.PrettyIndent
	}
	return "  "
}

/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"context"
	"strings"

	"github.com/blend/go-sdk/env"
)

// Config is the logger config.
type Config struct {
	// Flags hold the event types (i.e. flags) that are enabled.
	// If a flag is disabled, it is hidden from output _and_ listeners are not triggered.
	Flags	[]string	`json:"flags,omitempty" yaml:"flags,omitempty" env:"LOG_FLAGS,csv"`
	// Scopes hold the scope paths that are enabled.
	// If a scope is disabled, any events for that scope (or logger path) are hidden from output and listeners are not triggered.
	// A scope or path can be set on a logger with `sub := log.WithPath("foo", "bar")`.
	// It defaults to all scopes being enabled, or `*`.
	Scopes	[]string	`json:"scopes,omitempty" yaml:"scopes,omitempty" env:"LOG_SCOPES,csv"`
	// Writable holds event types (i.e. flags) that are shown in output.
	// If a flag is not writable, it is hidden from output but listeners _are_ triggered.
	// It defaults to all flags being writable, or `all`.
	Writable	[]string	`json:"writable,omitempty" yaml:"writable,omitempty" env:"LOG_WRITABLE,csv"`
	// WritableScopes are scopes that are shown in in output.
	// A scope can be set on a logger with `sub := log.WithPath("foo", "bar")`.
	// If a scope is not writable, it is hidden from output but listeners _are_ triggered.
	// It defaults to all scopes being writable, or `*`.
	WritableScopes	[]string	`json:"writableScopes,omitempty" yaml:"writableScopes,omitempty" env:"LOG_WRITABLE_SCOPES,csv"`
	// Format is the output format, either `text` or `json`.
	Format	string	`json:"format,omitempty" yaml:"format,omitempty" env:"LOG_FORMAT"`
	// Text holds text output specific options.
	Text	TextConfig	`json:"text,omitempty" yaml:"text,omitempty"`
	// JSON holds json specific options.
	JSON	JSONConfig	`json:"json,omitempty" yaml:"json,omitempty"`
}

// Resolve resolves the config.
func (c *Config) Resolve(ctx context.Context) error {
	return env.GetVars(ctx).ReadInto(c)
}

// FlagsOrDefault returns the enabled logger events.
func (c Config) FlagsOrDefault() []string {
	if len(c.Flags) > 0 {
		return c.Flags
	}
	return DefaultFlags
}

// WritableOrDefault returns the enabled logger events.
func (c Config) WritableOrDefault() []string {
	if len(c.Writable) > 0 {
		return c.Writable
	}
	return DefaultFlagsWritable
}

// ScopesOrDefault returns the enabled logger scopes.
func (c Config) ScopesOrDefault() []string {
	if len(c.Scopes) > 0 {
		return c.Scopes
	}
	return DefaultScopes
}

// WritableScopesOrDefault returns the writable logger scopes.
func (c Config) WritableScopesOrDefault() []string {
	if len(c.WritableScopes) > 0 {
		return c.WritableScopes
	}
	return DefaultWritableScopes
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
	HideTimestamp	bool	`json:"hideTimestamp,omitempty" yaml:"hideTimestamp,omitempty" env:"LOG_HIDE_TIMESTAMP"`
	HideFields	bool	`json:"hideFields,omitempty" yaml:"hideFields,omitempty" env:"LOG_HIDE_FIELDS"`
	NoColor		bool	`json:"noColor,omitempty" yaml:"noColor,omitempty" env:"NO_COLOR"`
	TimeFormat	string	`json:"timeFormat,omitempty" yaml:"timeFormat,omitempty" env:"LOG_TIME_FORMAT"`
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
	Pretty		bool	`json:"pretty,omitempty" yaml:"pretty,omitempty" env:"LOG_JSON_PRETTY"`
	PrettyPrefix	string	`json:"prettyPrefix,omitempty" yaml:"prettyPrefix,omitempty" env:"LOG_JSON_PRETTY_PREFIX"`
	PrettyIndent	string	`json:"prettyIndent,omitempty" yaml:"prettyIndent,omitempty" env:"LOG_JSON_PRETTY_INDENT"`
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

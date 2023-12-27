/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"io"

	"github.com/blend/go-sdk/env"
)

// Option is a logger option.
type Option func(*Logger) error

// OptConfig sets the logger based on a config.
func OptConfig(cfg Config) Option {
	return func(l *Logger) error {
		l.Formatter = cfg.Formatter()
		l.Flags = NewFlags(cfg.FlagsOrDefault()...)
		l.Writable = NewFlags(cfg.WritableOrDefault()...)
		l.Scopes = NewScopes(cfg.ScopesOrDefault()...)
		l.WritableScopes = NewScopes(cfg.WritableScopesOrDefault()...)
		return nil
	}
}

// OptConfigFromEnv sets the logger based on a config read from the environment.
// It will panic if there is an erro.
func OptConfigFromEnv() Option {
	return func(l *Logger) error {
		var cfg Config
		if err := env.Env().ReadInto(&cfg); err != nil {
			return err
		}
		l.Formatter = cfg.Formatter()
		l.Flags = NewFlags(cfg.FlagsOrDefault()...)
		l.Writable = NewFlags(cfg.WritableOrDefault()...)
		l.Scopes = NewScopes(cfg.ScopesOrDefault()...)
		l.WritableScopes = NewScopes(cfg.WritableScopesOrDefault()...)
		return nil
	}
}

/*
OptOutput sets the output writer for the logger.

It will wrap the output with a synchronizer if it's not already wrapped. You can also use this option to "unset" the output by passing in nil.

To set the output to be both stdout and a file, use the following:

	file, _ := os.Open("app.log")
	combined := io.MultiWriter(os.Stdout, file)
	log := logger.New(logger.OptOutput(combined))
*/
func OptOutput(output io.Writer) Option {
	return func(l *Logger) error {
		if output != nil {
			l.Output = NewInterlockedWriter(output)
		} else {
			l.Output = nil
		}
		return nil
	}
}

// OptPath sets an initial logger context path.
//
// This is useful if you want to label a logger to differentiate areas of an application
// but re-use the existing logger.
func OptPath(path ...string) Option {
	return func(l *Logger) error { l.Scope.Path = path; return nil }
}

// OptLabels sets an initial logger scope labels.
// This is useful if you want to add extra information for events by default (like environment).
func OptLabels(labels ...Labels) Option {
	return func(l *Logger) error { l.Scope.Labels = CombineLabels(labels...); return nil }
}

// OptJSON sets the output formatter for the logger as json.
func OptJSON(opts ...JSONOutputFormatterOption) Option {
	return func(l *Logger) error { l.Formatter = NewJSONOutputFormatter(opts...); return nil }
}

// OptText sets the output formatter for the logger as json.
func OptText(opts ...TextOutputFormatterOption) Option {
	return func(l *Logger) error { l.Formatter = NewTextOutputFormatter(opts...); return nil }
}

// OptFormatter sets the output formatter.
func OptFormatter(formatter WriteFormatter) Option {
	return func(l *Logger) error { l.Formatter = formatter; return nil }
}

// OptFlags sets the flags on the logger.
func OptFlags(flags *Flags) Option {
	return func(l *Logger) error { l.Flags = flags; return nil }
}

// OptWritable sets the writable flags on the logger.
func OptWritable(flags *Flags) Option {
	return func(l *Logger) error { l.Writable = flags; return nil }
}

// OptScopes sets the scopes on the logger.
func OptScopes(scopes *Scopes) Option {
	return func(l *Logger) error { l.Scopes = scopes; return nil }
}

// OptWritableScopes sets the writable scopes on the logger.
func OptWritableScopes(scopes *Scopes) Option {
	return func(l *Logger) error { l.WritableScopes = scopes; return nil }
}

// OptAll sets all flags enabled on the logger by default.
func OptAll() Option {
	return func(l *Logger) error { l.Flags.SetAll(); return nil }
}

// OptAllWritable sets all flags enabled on the logger by default.
func OptAllWritable() Option {
	return func(l *Logger) error { l.Writable.SetAll(); return nil }
}

// OptAllScopes sets all scopes enabled on the logger by default.
func OptAllScopes() Option {
	return func(l *Logger) error { l.Scopes.SetAll(); return nil }
}

// OptAllWritableScopes sets all scopes for writing enabled on the logger by default.
func OptAllWritableScopes() Option {
	return func(l *Logger) error { l.WritableScopes.SetAll(); return nil }
}

// OptNone sets no flags enabled on the logger by default.
func OptNone() Option {
	return func(l *Logger) error { l.Flags.SetNone(); return nil }
}

// OptNoneWritable sets no flags enabled for writing on the logger by default.
func OptNoneWritable() Option {
	return func(l *Logger) error { l.Writable.SetNone(); return nil }
}

// OptNoneScopes sets no scopes enabled on the logger by default.
func OptNoneScopes() Option {
	return func(l *Logger) error { l.Scopes.SetNone(); return nil }
}

// OptNoneWritableScopes sets no scopes enabled for writing on the logger by default.
func OptNoneWritableScopes() Option {
	return func(l *Logger) error { l.WritableScopes.SetNone(); return nil }
}

// OptEnabled sets enabled flags on the logger.
func OptEnabled(flags ...string) Option {
	return func(l *Logger) error { l.Flags.Enable(flags...); return nil }
}

// OptEnabledWritable sets enabled writable flags on the logger.
func OptEnabledWritable(flags ...string) Option {
	return func(l *Logger) error { l.Writable.Enable(flags...); return nil }
}

// OptEnabledScopes sets enabled scopes on the logger.
func OptEnabledScopes(scopes ...string) Option {
	return func(l *Logger) error { l.Scopes.Enable(scopes...); return nil }
}

// OptEnabledWritableScopes sets enabled writable scopes on the logger.
func OptEnabledWritableScopes(scopes ...string) Option {
	return func(l *Logger) error { l.WritableScopes.Enable(scopes...); return nil }
}

// OptDisabled sets disabled flags on the logger.
func OptDisabled(flags ...string) Option {
	return func(l *Logger) error { l.Flags.Disable(flags...); return nil }
}

// OptDisabledWritable sets disabled flags on the logger.
func OptDisabledWritable(flags ...string) Option {
	return func(l *Logger) error { l.Writable.Disable(flags...); return nil }
}

// OptDisabledScopes sets disabled scopes on the logger.
func OptDisabledScopes(scopes ...string) Option {
	return func(l *Logger) error { l.Scopes.Disable(scopes...); return nil }
}

// OptDisabledWritableScopes sets disabled flags on the logger.
func OptDisabledWritableScopes(scopes ...string) Option {
	return func(l *Logger) error { l.WritableScopes.Disable(scopes...); return nil }
}

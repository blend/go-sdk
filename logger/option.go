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

// OptSubContext sets an initial logger context path.
// This is useful if you want to label a logger to differentiate multiple loggers.
func OptSubContext(path ...string) Option {
	return func(l *Logger) error { l.Context.Path = path; return nil }
}

// OptFields sets an initial logger context fields.
func OptFields(fields Fields) Option {
	return func(l *Logger) error { l.Context.Fields = fields; return nil }
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

// OptAll sets all flags enabled on the logger by default.
func OptAll() Option {
	return func(l *Logger) error { l.Flags.SetAll(); return nil }
}

// OptNone sets no flags enabled on the logger by default.
func OptNone() Option {
	return func(l *Logger) error { l.Flags.SetNone(); return nil }
}

// OptEnabled sets enabled flags on the logger.
func OptEnabled(flags ...string) Option {
	return func(l *Logger) error { l.Flags.Enable(flags...); return nil }
}

// OptDisabled sets disabled flags on the logger.
func OptDisabled(flags ...string) Option {
	return func(l *Logger) error { l.Flags.Disable(flags...); return nil }
}

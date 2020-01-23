package graceful

import (
	"os"
	"time"
)

// OptShutdownSignal sets the shutdown signal.
func OptShutdownSignal(signal chan os.Signal) Option {
	return func(so *Options) { so.ShutdownSignal = signal }
}

// OptUpdateSignal sets the update signal.
func OptUpdateSignal(signal chan os.Signal) Option {
	return func(so *Options) { so.UpdateSignal = signal }
}

// OptShutdownGracePeriod sets a shutdown grace period.
func OptShutdownGracePeriod(gracePeriod time.Duration) Option {
	return func(so *Options) { so.ShutdownGracePeriod = gracePeriod }
}

// OptUpdateErrors sets an error collector channel for errors.
func OptUpdateErrors(errs chan error) Option {
	return func(so *Options) { so.UpdateErrors = errs }
}

// Option is a mutator for shutdown options.
type Option func(*Options)

// Options are the options for graceful shutdown.
type Options struct {
	ShutdownGracePeriod time.Duration
	ShutdownSignal      chan os.Signal
	UpdateSignal        chan os.Signal
	UpdateErrors        chan error
}

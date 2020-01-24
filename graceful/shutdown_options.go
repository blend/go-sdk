package graceful

import (
	"os"
	"time"
)

// OptDefaultShutdownSignal returns an option that sets the shutdown signal to the defaults.
func OptDefaultShutdownSignal() ShutdownOption {
	return func(so *ShutdownOptions) { so.ShutdownSignal = Notify(DefaultShutdownSignals...) }
}

// OptShutdownSignal sets the shutdown signal.
func OptShutdownSignal(signal chan os.Signal) ShutdownOption {
	return func(so *ShutdownOptions) { so.ShutdownSignal = signal }
}

// OptShutdownGracePeriod sets a shutdown grace period.
func OptShutdownGracePeriod(gracePeriod time.Duration) ShutdownOption {
	return func(so *ShutdownOptions) { so.ShutdownGracePeriod = gracePeriod }
}

// ShutdownOption is a mutator for shutdown options.
type ShutdownOption func(*ShutdownOptions)

// ShutdownOptions are the options for graceful shutdown.
type ShutdownOptions struct {
	ShutdownGracePeriod time.Duration
	ShutdownSignal      chan os.Signal
}

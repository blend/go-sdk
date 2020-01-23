package graceful

import (
	"os"
	"syscall"
)

// DefaultShutdownSignals are the default os signals to capture to shut down.
var DefaultShutdownSignals = []os.Signal{
	os.Interrupt, syscall.SIGTERM,
}

// DefaultUpdateSignals are the default os signals to capture for updates.
var DefaultUpdateSignals = []os.Signal{
	syscall.SIGHUP,
}

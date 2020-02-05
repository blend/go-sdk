package graceful

import (
	"os"
	"os/signal"
)

// Notify returns a channel that listens for a given set of os signals.
func Notify(signals ...os.Signal) chan os.Signal {
	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, signals...)
	return terminateSignal
}

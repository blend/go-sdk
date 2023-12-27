/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package graceful

import (
	"os"
	"os/signal"
)

// Notify returns a channel that listens for a given set of os signals.
func Notify(signals ...os.Signal) chan os.Signal {
	return NotifyWithCapacity(1, signals...)
}

// NotifyWithCapacity returns a channel with a given capacity
// that listens for a given set of os signals.
func NotifyWithCapacity(capacity int, signals ...os.Signal) chan os.Signal {
	terminateSignal := make(chan os.Signal, capacity)
	signal.Notify(terminateSignal, signals...)
	return terminateSignal
}

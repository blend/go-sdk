/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package graceful

import (
	"os"
	"syscall"
)

// DefaultShutdownSignals are the default os signals to capture to shut down.
var DefaultShutdownSignals = []os.Signal{
	os.Interrupt, syscall.SIGTERM,
}

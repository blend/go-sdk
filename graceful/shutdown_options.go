/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package graceful

import (
	"os"
)

// OptDefaultShutdownSignal returns an option that sets the shutdown signal to the defaults.
func OptDefaultShutdownSignal() ShutdownOption {
	return func(so *ShutdownOptions) { so.ShutdownSignal = Notify(DefaultShutdownSignals...) }
}

// OptShutdownSignal sets the shutdown signal.
func OptShutdownSignal(signal chan os.Signal) ShutdownOption {
	return func(so *ShutdownOptions) { so.ShutdownSignal = signal }
}

// ShutdownOption is a mutator for shutdown options.
type ShutdownOption func(*ShutdownOptions)

// ShutdownOptions are the options for graceful shutdown.
type ShutdownOptions struct {
	ShutdownSignal chan os.Signal
}

/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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

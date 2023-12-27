/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package logger

import (
	stdlog "log"
)

// StdlibShim returns a stdlib logger that writes to a given logger instance.
func StdlibShim(log Triggerable, opts ...ShimWriterOption) *stdlog.Logger {
	shim := NewShimWriter(log)
	for _, opt := range opts {
		opt(&shim)
	}
	return stdlog.New(shim, "", 0)
}

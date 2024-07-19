/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package envoyutil_test

import (
	"bytes"

	"github.com/blend/go-sdk/logger"
)

// InMemoryLog creates a logger that logs to the in-memory buffer passed in.
func InMemoryLog(logBuffer *bytes.Buffer) logger.Log {
	return logger.MustNew(
		logger.OptAll(),
		logger.OptOutput(logBuffer),
		logger.OptFormatter(logger.NewTextOutputFormatter(
			logger.OptTextNoColor(),
			logger.OptTextHideTimestamp(),
		)),
	)
}

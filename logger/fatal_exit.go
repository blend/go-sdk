/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"os"
	"sync"
)

var (
	_log		*Logger
	_logInit	sync.Once
)

func ensureLog() {
	_logInit.Do(func() { _log = MustNew(OptEnabled(Info, Debug, Warning, Error, Fatal)) })
}

// FatalExit will print the error and exit the process with exit(1).
func FatalExit(err error) {
	ensureLog()
	_log.Fatal(err)
	os.Exit(1)
}

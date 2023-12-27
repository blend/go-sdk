/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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

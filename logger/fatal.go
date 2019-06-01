package logger

import (
	"os"
	"sync"
)

var (
	_log     *Logger
	_logInit sync.Once
)

func ensureLog() {
	_logInit.Do(func() { _log = MustNew(OptEnabled(Info, Debug, Warning, Error, Fatal)) })
}

// MaybeFatalExit will print the error and exit the process
// with exit(1) if the error isn't nil.
func MaybeFatalExit(err error) {
	if err == nil {
		return
	}
	FatalExit(err)
}

// FatalExit will print the error and exit the process with exit(1).
func FatalExit(err error) {
	ensureLog()
	_log.Fatal(err)
	os.Exit(1)
}

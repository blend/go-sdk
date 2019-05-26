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

// SubContext returns a new default sub context.
func SubContext(heading string, options ...ContextOption) Context {
	ensureLog()
	return _log.SubContext(heading, options...)
}

// Infof prints an info message with the default logger.
func Infof(format string, args ...interface{}) {
	ensureLog()
	_log.Infof(format, args...)
}

// Debugf prints an debug message with the default logger.
func Debugf(format string, args ...interface{}) {
	ensureLog()
	_log.Debugf(format, args...)
}

// Warningf prints an warning message with the default logger.
func Warningf(format string, args ...interface{}) {
	ensureLog()
	_log.Warningf(format, args...)
}

// Errorf prints an error message with the default logger.
func Errorf(format string, args ...interface{}) {
	ensureLog()
	_log.Errorf(format, args...)
}

// Fatalf prints an fatal message with the default logger.
func Fatalf(format string, args ...interface{}) {
	ensureLog()
	_log.Fatalf(format, args...)
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

package logger

import (
	"os"
	"sync"
)

var (
	log     *Logger
	logInit sync.Once
)

func ensureLog() {
	logInit.Do(func() { log = MustNew() })
}

// SubContext returns a new default sub context.
func SubContext(heading string, options ...ContextOption) Context {
	ensureLog()
	return log.SubContext(heading, options...)
}

// Infof prints an info message with the default logger.
func Infof(format string, args ...interface{}) {
	ensureLog()
	log.Infof(format, args...)
}

// Debugf prints an debug message with the default logger.
func Debugf(format string, args ...interface{}) {
	ensureLog()
	log.Debugf(format, args...)
}

// Warningf prints an warning message with the default logger.
func Warningf(format string, args ...interface{}) {
	ensureLog()
	log.Warningf(format, args...)
}

// Errorf prints an error message with the default logger.
func Errorf(format string, args ...interface{}) {
	ensureLog()
	log.Errorf(format, args...)
}

// Fatalf prints an fatal message with the default logger.
func Fatalf(format string, args ...interface{}) {
	ensureLog()
	log.Errorf(format, args...)
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
	log.Fatal(err)
	os.Exit(1)
}

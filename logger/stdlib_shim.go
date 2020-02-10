package logger

import (
	"context"
	"log"
)

// StdlibShim returns a stdlib logger that writes to a given logger instance.
func StdlibShim(ctx context.Context, flag string, handler Triggerable) *log.Logger {
	return log.New(ShimWriter{Context: ctx, Flag: flag, Log: handler}, "", 0)
}

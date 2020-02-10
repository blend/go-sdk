package loggerutil

import (
	"context"
	"log"

	"github.com/blend/go-sdk/logger"
)

// StdlibShim returns a stdlib logger that writes to a given logger instance.
func StdlibShim(ctx context.Context, flag string, handler logger.Triggerable) *log.Logger {
	return log.New(logger.ShimWriter{Context: ctx, Flag: flag, Log: handler}, "", log.LstdFlags)
}

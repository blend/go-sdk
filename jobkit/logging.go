package jobkit

import (
	"context"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/logger"
)

// LogScoperDebugf is a logger interface.
type LogScoperDebugf interface {
	logger.Scoper
	logger.DebugfReceiver
}

// LogScoperInfof is a logger interface.
type LogScoperInfof interface {
	logger.Scoper
	logger.InfofReceiver
}

// LogScoperWarningf is a logger interface.
type LogScoperWarningf interface {
	logger.Scoper
	logger.WarningfReceiver
}

// LogScoperWarning is a logger interface.
type LogScoperWarning interface {
	logger.Scoper
	logger.WarningReceiver
}

// LogScoperErrorf is a logger interface.
type LogScoperErrorf interface {
	logger.Scoper
	logger.ErrorfReceiver
}

// LogScoperError is a logger interface.
type LogScoperError interface {
	logger.Scoper
	logger.ErrorReceiver
}

// LogScoperFatalf is a logger interface.
type LogScoperFatalf interface {
	logger.Scoper
	logger.FatalfReceiver
}

// LogScoperFatal is a logger interface.
type LogScoperFatal interface {
	logger.Scoper
	logger.FatalReceiver
}

// Debugf prints an info message if the logger is set.
func Debugf(ctx context.Context, log LogScoperDebugf, format string, args ...interface{}) {
	if log == nil {
		return
	}
	ji := cron.GetJobInvocation(ctx)
	log.WithPath(ji.ID).Debugf(format, args...)
}

// Infof prints an info message if the logger is set.
func Infof(ctx context.Context, log LogScoperInfof, format string, args ...interface{}) {
	if log == nil {
		return
	}
	ji := cron.GetJobInvocation(ctx)
	log.WithPath(ji.ID).Infof(format, args...)
}

// Warningf prints a warning message if the logger is set.
func Warningf(ctx context.Context, log LogScoperWarningf, format string, args ...interface{}) {
	if log == nil {
		return
	}
	ji := cron.GetJobInvocation(ctx)
	log.WithPath(ji.ID).Warningf(format, args...)
}

// Warning prints an warning if the logger is set.
func Warning(ctx context.Context, log LogScoperWarning, err error) {
	if log == nil || err == nil {
		return
	}
	ji := cron.GetJobInvocation(ctx)
	log.WithPath(ji.ID).Warning(err)
}

// Errorf prints an error message if the logger is set.
func Errorf(ctx context.Context, log LogScoperErrorf, format string, args ...interface{}) {
	if log == nil {
		return
	}
	ji := cron.GetJobInvocation(ctx)
	log.WithPath(ji.ID).Errorf(format, args...)
}

// Error prints an error if the logger is set.
func Error(ctx context.Context, log LogScoperError, err error) {
	if log == nil || err == nil {
		return
	}
	ji := cron.GetJobInvocation(ctx)
	log.WithPath(ji.ID).Error(err)
}

// Fatalf prints a fatal error message if the logger is set.
func Fatalf(ctx context.Context, log LogScoperFatalf, format string, args ...interface{}) {
	if log == nil {
		return
	}
	ji := cron.GetJobInvocation(ctx)
	log.WithPath(ji.ID).Fatalf(format, args...)
}

// Fatal prints a fatal error if the logger is set.
func Fatal(ctx context.Context, log LogScoperFatal, err error) {
	if log == nil || err == nil {
		return
	}
	ji := cron.GetJobInvocation(ctx)
	log.WithPath(ji.ID).Fatal(err)
}

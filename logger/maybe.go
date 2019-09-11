package logger

import (
	"context"
)

// IsLoggerSet returns if the logger instance is set.
func IsLoggerSet(log interface{}) bool {
	if log == nil {
		return false
	}
	if typed, ok := log.(*Logger); ok {
		return typed != nil
	}
	return true
}

// MaybeTrigger triggers an event if the logger is set.
func MaybeTrigger(ctx context.Context, log Triggerable, e Event) {
	if !IsLoggerSet(log) {
		return
	}
	log.Trigger(ctx, e)
}

// MaybeInfo triggers Info if the logger is set.
func MaybeInfo(log InfoReceiver, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.Info(args...)
}

// MaybeInfoContext triggers Info in a given context if the logger.
func MaybeInfoContext(ctx context.Context, log Scoper, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.WithContext(ctx).Info(args...)
}

// MaybeInfof triggers Infof if the logger is set.
func MaybeInfof(log InfofReceiver, format string, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.Infof(format, args...)
}

// MaybeInfofContext triggers Infof in a given context if the logger is set.
func MaybeInfofContext(ctx context.Context, log Scoper, format string, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.WithContext(ctx).Infof(format, args...)
}

// MaybeDebug triggers Debug if the logger is set.
func MaybeDebug(log DebugReceiver, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.Debug(args...)
}

// MaybeDebugContext triggers Debug in a given context if the logger is set.
func MaybeDebugContext(ctx context.Context, log Scoper, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.WithContext(ctx).Debug(args...)
}

// MaybeDebugf triggers Debugf if the logger is set.
func MaybeDebugf(log DebugfReceiver, format string, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.Debugf(format, args...)
}

// MaybeDebugfContext triggers Debugf in a given context if the logger is set.
func MaybeDebugfContext(ctx context.Context, log Scoper, format string, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.WithContext(ctx).Debugf(format, args...)
}

// MaybeWarningf triggers Warningf if the logger is set.
func MaybeWarningf(log WarningfReceiver, format string, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.Warningf(format, args...)
}

// MaybeWarningfContext triggers Warningf in a given context if the logger is set.
func MaybeWarningfContext(ctx context.Context, log Scoper, format string, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.WithContext(ctx).Warningf(format, args...)
}

// MaybeWarning triggers Warning if the logger is set.
func MaybeWarning(log WarningReceiver, err error) {
	if !IsLoggerSet(log) || err == nil {
		return
	}
	log.Warning(err)
}

// MaybeWarningContext triggers Warning in a given context if the logger is set.
func MaybeWarningContext(ctx context.Context, log Scoper, err error) {
	if !IsLoggerSet(log) || err == nil {
		return
	}
	log.WithContext(ctx).Warning(err)
}

// MaybeErrorf triggers Errorf if the logger is set.
func MaybeErrorf(log ErrorfReceiver, format string, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.Errorf(format, args...)
}

// MaybeErrorfContext triggers Errorf in a given context if the logger is set.
func MaybeErrorfContext(ctx context.Context, log Scoper, format string, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.WithContext(ctx).Errorf(format, args...)
}

// MaybeError triggers Error if the logger is set.
func MaybeError(log ErrorReceiver, err error) {
	if !IsLoggerSet(log) || err == nil {
		return
	}
	log.Error(err)
}

// MaybeErrorContext triggers Error in a given context if the logger is set.
func MaybeErrorContext(ctx context.Context, log Scoper, err error) {
	if !IsLoggerSet(log) || err == nil {
		return
	}
	log.WithContext(ctx).Error(err)
}

// MaybeFatalf triggers Fatalf if the logger is set.
func MaybeFatalf(log FatalfReceiver, format string, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.Fatalf(format, args...)
}

// MaybeFatalfContext triggers Fatalf in a given context if the logger is set.
func MaybeFatalfContext(ctx context.Context, log Scoper, format string, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.WithContext(ctx).Fatalf(format, args...)
}

// MaybeFatal triggers Fatal if the logger is set.
func MaybeFatal(log FatalReceiver, err error) {
	if !IsLoggerSet(log) || err == nil {
		return
	}
	log.Fatal(err)
}

// MaybeFatalContext triggers Fatal in a given context if the logger is set.
func MaybeFatalContext(ctx context.Context, log Scoper, err error) {
	if !IsLoggerSet(log) || err == nil {
		return
	}
	log.WithContext(ctx).Fatal(err)
}

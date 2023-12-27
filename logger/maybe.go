/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package logger

import (
	"context"
	"os"
)

// IsLoggerSet returns if the logger instance is set.
func IsLoggerSet(log interface{}) bool {
	if typed, ok := log.(*Logger); ok {
		return typed != nil
	}
	return log != nil
}

// MaybeDrainContext drains a logger if it's a valid reference and can be drained.
func MaybeDrainContext(ctx context.Context, log interface{}) {
	if typed, ok := log.(*Logger); ok && typed != nil {
		typed.DrainContext(ctx)
	}
}

// MaybeTrigger triggers an event if the logger is set.
//
// DEPRECATION(1.2021*): this method will be changed to drop the context and use `context.Background()`.
func MaybeTrigger(ctx context.Context, log Triggerable, e Event) {
	if !IsLoggerSet(log) {
		return
	}
	log.TriggerContext(ctx, e)
}

// MaybeTriggerContext triggers an event if the logger is set in a given context.
func MaybeTriggerContext(ctx context.Context, log Triggerable, e Event) {
	if !IsLoggerSet(log) {
		return
	}
	log.TriggerContext(ctx, e)
}

// MaybeInfo triggers Info if the logger is set.
func MaybeInfo(log InfoReceiver, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.Info(args...)
}

// MaybeInfoContext triggers Info in a given context if the logger.
func MaybeInfoContext(ctx context.Context, log InfoReceiver, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.InfoContext(ctx, args...)
}

// MaybeInfof triggers Infof if the logger is set.
func MaybeInfof(log InfofReceiver, format string, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.Infof(format, args...)
}

// MaybeInfofContext triggers Infof in a given context if the logger is set.
func MaybeInfofContext(ctx context.Context, log InfofReceiver, format string, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.InfofContext(ctx, format, args...)
}

// MaybeDebug triggers Debug if the logger is set.
func MaybeDebug(log DebugReceiver, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.Debug(args...)
}

// MaybeDebugContext triggers Debug in a given context if the logger is set.
func MaybeDebugContext(ctx context.Context, log DebugReceiver, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.DebugContext(ctx, args...)
}

// MaybeDebugf triggers Debugf if the logger is set.
func MaybeDebugf(log DebugfReceiver, format string, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.Debugf(format, args...)
}

// MaybeDebugfContext triggers Debugf in a given context if the logger is set.
func MaybeDebugfContext(ctx context.Context, log DebugfReceiver, format string, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.DebugfContext(ctx, format, args...)
}

// MaybeWarningf triggers Warningf if the logger is set.
func MaybeWarningf(log WarningfReceiver, format string, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.Warningf(format, args...)
}

// MaybeWarningfContext triggers Warningf in a given context if the logger is set.
func MaybeWarningfContext(ctx context.Context, log WarningfReceiver, format string, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.WarningfContext(ctx, format, args...)
}

// MaybeWarning triggers Warning if the logger is set.
func MaybeWarning(log WarningReceiver, err error) {
	if !IsLoggerSet(log) || err == nil {
		return
	}
	log.Warning(err)
}

// MaybeWarningContext triggers Warning in a given context if the logger is set.
func MaybeWarningContext(ctx context.Context, log WarningReceiver, err error) {
	if !IsLoggerSet(log) || err == nil {
		return
	}
	log.WarningContext(ctx, err)
}

// MaybeErrorf triggers Errorf if the logger is set.
func MaybeErrorf(log ErrorfReceiver, format string, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.Errorf(format, args...)
}

// MaybeErrorfContext triggers Errorf in a given context if the logger is set.
func MaybeErrorfContext(ctx context.Context, log ErrorfReceiver, format string, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.ErrorfContext(ctx, format, args...)
}

// MaybeError triggers Error if the logger is set.
func MaybeError(log ErrorReceiver, err error) {
	if !IsLoggerSet(log) || err == nil {
		return
	}
	log.Error(err)
}

// MaybeErrorContext triggers Error in a given context if the logger is set.
func MaybeErrorContext(ctx context.Context, log ErrorReceiver, err error) {
	if !IsLoggerSet(log) || err == nil {
		return
	}
	log.ErrorContext(ctx, err)
}

// MaybeFatalf triggers Fatalf if the logger is set.
func MaybeFatalf(log FatalfReceiver, format string, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.Fatalf(format, args...)
}

// MaybeFatalfContext triggers Fatalf in a given context if the logger is set.
func MaybeFatalfContext(ctx context.Context, log FatalfReceiver, format string, args ...interface{}) {
	if !IsLoggerSet(log) {
		return
	}
	log.FatalfContext(ctx, format, args...)
}

// MaybeFatal triggers Fatal if the logger is set.
func MaybeFatal(log FatalReceiver, err error) {
	if !IsLoggerSet(log) || err == nil {
		return
	}
	log.Fatal(err)
}

// MaybeFatalContext triggers Fatal in a given context if the logger is set.
func MaybeFatalContext(ctx context.Context, log FatalReceiver, err error) {
	if !IsLoggerSet(log) || err == nil {
		return
	}
	log.FatalContext(ctx, err)
}

// MaybeFatalExit triggers Fatal if the logger is set and the error is set, and exit(1)s.
func MaybeFatalExit(log FatalCloser, err error) {
	if !IsLoggerSet(log) || err == nil {
		return
	}
	log.Fatal(err)
	log.Close()
	os.Exit(1)
}

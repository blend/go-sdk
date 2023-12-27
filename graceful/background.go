/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package graceful

import (
	"context"
	"os"
	"os/signal"
)

// OptBackgroundSignals sets the signals.
func OptBackgroundSignals(signals ...os.Signal) BackgroundOption {
	return func(bo *BackgroundOptions) { bo.Signals = signals }
}

// OptBackgroundLog sets the logger.
func OptBackgroundLog(log Logger) BackgroundOption {
	return func(bo *BackgroundOptions) { bo.Log = log }
}

// OptBackgroundSkipStopOnSignal sets if we should stop the signal channel on stop.
func OptBackgroundSkipStopOnSignal(skipStopOnSignal bool) BackgroundOption {
	return func(bo *BackgroundOptions) { bo.SkipStopOnSignal = skipStopOnSignal }
}

// BackgroundOption mutates background options
type BackgroundOption func(*BackgroundOptions)

// BackgroundOptions are options for the background context.
type BackgroundOptions struct {
	Signals          []os.Signal
	Log              Logger
	SkipStopOnSignal bool
}

// Background yields a context that will signal `<-ctx.Done()` when
// a signal is sent to the process (as specified in `DefaultShutdownSignals`).
func Background(opts ...BackgroundOption) context.Context {
	options := BackgroundOptions{
		Signals:          DefaultShutdownSignals,
		SkipStopOnSignal: false,
	}
	for _, opt := range opts {
		opt(&options)
	}

	ctx, cancel := context.WithCancel(context.Background())
	shutdown := Notify(options.Signals...)

	go func() {
		MaybeDebugf(options.Log, "graceful background; waiting for shutdown signal")
		<-shutdown

		MaybeDebugf(options.Log, "graceful background; shutdown signal received, canceling context")
		cancel()

		if !options.SkipStopOnSignal {
			MaybeDebugf(options.Log, "graceful background; stopping shutdown signal channel")
			signal.Stop(shutdown)
		}
	}()
	return ctx
}

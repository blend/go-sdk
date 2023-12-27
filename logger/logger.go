/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"context"
	"io"
	"os"
	"sync"
)

// New returns a new logger with a given set of enabled flags.
// By default it uses a text output formatter writing to stdout.
func New(options ...Option) (*Logger, error) {
	l := &Logger{
		Formatter:      NewTextOutputFormatter(),
		Output:         NopCloserWriter{NewInterlockedWriter(os.Stdout)},
		RecoverPanics:  DefaultRecoverPanics,
		Flags:          NewFlags(DefaultFlags...),
		Writable:       FlagsAll(),
		Scopes:         ScopesAll(),
		WritableScopes: ScopesAll(),
	}

	l.Scope = NewScope(l)
	var err error
	for _, option := range options {
		if err = option(l); err != nil {
			return nil, err
		}
	}
	return l, nil
}

// MustNew creates a new logger with a given list of options and panics on error.
func MustNew(options ...Option) *Logger {
	log, err := New(options...)
	if err != nil {
		panic(err)
	}
	return log
}

// All returns a new logger with all flags enabled.
func All(options ...Option) *Logger {
	return MustNew(
		append([]Option{
			OptConfigFromEnv(),
			OptAll(),
		}, options...)...)
}

// None returns a new logger with all flags enabled.
func None(options ...Option) *Logger {
	return MustNew(
		append([]Option{
			OptNone(),
			OptOutput(nil),
			OptFormatter(nil),
		}, options...)...)
}

// Prod returns a new logger tuned for production use.
// It writes to os.Stderr with text output colorization disabled.
func Prod(options ...Option) *Logger {
	return MustNew(
		append([]Option{
			OptAll(),
			OptOutput(os.Stderr),
			OptFormatter(NewTextOutputFormatter(OptTextNoColor())),
		}, options...)...)
}

// Memory creates a logger that logs to the in-memory writer passed in.
//
// It is useful for writing tests that collect log output.
func Memory(buffer io.Writer, options ...Option) *Logger {
	return MustNew(
		append([]Option{
			OptAll(),
			OptOutput(buffer),
			OptFormatter(NewTextOutputFormatter(
				OptTextNoColor(),
				OptTextHideTimestamp(),
			)),
		}, options...)...)
}

// Logger is a handler for various logging events with descendent handlers.
type Logger struct {
	sync.Mutex
	*Flags
	Scope

	Writable       *Flags
	Scopes         *Scopes
	WritableScopes *Scopes
	RecoverPanics  bool

	Output    io.Writer
	Formatter WriteFormatter
	Errors    chan error

	// Filters hold filters organized by flag, and then by filter name.
	// The intent is to modify event data before it is written or given to listeners.
	Filters map[string]map[string]Filter
	// Listeners hold event listeners organized by flag, and then by listener name.
	Listeners map[string]map[string]*Worker
}

// GetFlags returns the flags.
func (l *Logger) GetFlags() *Flags {
	return l.Flags
}

// GetWritable returns the writable flags.
func (l *Logger) GetWritable() *Flags {
	return l.Writable
}

// HasListeners returns if there are registered listener for an event.
func (l *Logger) HasListeners(flag string) bool {
	l.Lock()
	defer l.Unlock()

	if l.Listeners == nil {
		return false
	}
	listeners, ok := l.Listeners[flag]
	if !ok {
		return false
	}
	return len(listeners) > 0
}

// HasListener returns if a specific listener is registered for a flag.
func (l *Logger) HasListener(flag, listenerName string) bool {
	l.Lock()
	defer l.Unlock()

	if l.Listeners == nil {
		return false
	}
	workers, ok := l.Listeners[flag]
	if !ok {
		return false
	}
	_, ok = workers[listenerName]
	return ok
}

// Listen adds a listener for a given flag.
func (l *Logger) Listen(flag, listenerName string, listener Listener) {
	l.Lock()
	defer l.Unlock()

	if l.Listeners == nil {
		l.Listeners = make(map[string]map[string]*Worker)
	}
	if l.Listeners[flag] == nil {
		l.Listeners[flag] = make(map[string]*Worker)
	}

	eventListener := NewWorker(listener)
	l.Listeners[flag][listenerName] = eventListener
	go func() { _ = eventListener.Start() }()
	<-eventListener.NotifyStarted()
}

// RemoveListeners clears *all* listeners for a Flag.
func (l *Logger) RemoveListeners(flag string) error {
	l.Lock()
	defer l.Unlock()

	if l.Listeners == nil {
		return nil
	}

	listeners, ok := l.Listeners[flag]
	if !ok {
		return nil
	}
	var err error
	for _, l := range listeners {
		if err = l.Stop(); err != nil {
			return err
		}
	}
	delete(l.Listeners, flag)
	return nil
}

// RemoveListener clears a specific listener for a Flag.
func (l *Logger) RemoveListener(flag, listenerName string) error {
	l.Lock()
	defer l.Unlock()

	if l.Listeners == nil {
		return nil
	}

	listeners, ok := l.Listeners[flag]
	if !ok {
		return nil
	}

	worker, ok := listeners[listenerName]
	if !ok {
		return nil
	}
	if err := worker.Stop(); err != nil {
		return err
	}
	delete(listeners, listenerName)
	if len(listeners) == 0 {
		delete(l.Listeners, flag)
	}
	return nil
}

// HasFilters returns if a logger has filters for a given flag.
func (l *Logger) HasFilters(flag string) bool {
	l.Lock()
	defer l.Unlock()

	if l.Filters == nil {
		return false
	}
	filters, ok := l.Filters[flag]
	if !ok {
		return false
	}
	return len(filters) > 0
}

// HasFilter returns if a logger has a given filter by name.
func (l *Logger) HasFilter(flag, filterName string) bool {
	l.Lock()
	defer l.Unlock()

	if l.Filters == nil {
		return false
	}
	filters, ok := l.Filters[flag]
	if !ok {
		return false
	}
	_, ok = filters[filterName]
	return ok
}

// Filter adds a given filter for a given flag.
func (l *Logger) Filter(flag, filterName string, filter Filter) {
	l.Lock()
	defer l.Unlock()

	if l.Filters == nil {
		l.Filters = make(map[string]map[string]Filter)
	}
	if l.Filters[flag] == nil {
		l.Filters[flag] = make(map[string]Filter)
	}
	l.Filters[flag][filterName] = filter
}

// RemoveFilters clears *all* filters for a Flag.
func (l *Logger) RemoveFilters(flag string) {
	l.Lock()
	defer l.Unlock()
	delete(l.Filters, flag)
}

// RemoveFilter clears a specific filter for a Flag.
func (l *Logger) RemoveFilter(flag, filterName string) {
	l.Lock()
	defer l.Unlock()

	if l.Filters == nil {
		return
	}
	filters, ok := l.Filters[flag]
	if !ok {
		return
	}
	delete(filters, filterName)
}

// Dispatch fires the listeners for a given event asynchronously, and writes the event to the output.
// The invocations will be queued in a work queue per listener.
// There are no order guarantees on when these events will be processed across listeners.
// This call will not block on the event listeners, but will block on the write.
func (l *Logger) Dispatch(ctx context.Context, e Event) {
	if e == nil {
		return
	}
	flag := e.GetFlag()
	if !l.IsEnabled(flag) {
		return
	}
	if !l.Scopes.IsEnabled(GetPath(ctx)...) {
		return
	}

	if !IsSkipTrigger(ctx) {
		var filters map[string]Filter
		var listeners map[string]*Worker
		l.Lock()
		if l.Filters != nil {
			if flagFilters, ok := l.Filters[flag]; ok {
				filters = flagFilters
			}
		}
		if l.Listeners != nil {
			if flagListeners, ok := l.Listeners[flag]; ok {
				listeners = flagListeners
			}
		}
		l.Unlock()

		var shouldFilter bool
		for _, filter := range filters {
			e, shouldFilter = filter(ctx, e)
			if shouldFilter {
				return
			}
		}
		for _, listener := range listeners {
			listener.Work <- EventWithContext{ctx, e}
		}
	}

	l.Write(ctx, e)
}

// Write writes an event synchronously to the writer either as a normal even or as an error.
func (l *Logger) Write(ctx context.Context, e Event) {
	// if a formater or the output are unset, bail.
	if l.Formatter == nil || l.Output == nil {
		return
	}

	if IsSkipWrite(ctx) {
		return
	}
	if !l.Writable.IsEnabled(e.GetFlag()) {
		return
	}
	if !l.WritableScopes.IsEnabled(GetPath(ctx)...) {
		return
	}

	err := l.Formatter.WriteFormat(ctx, l.Output, e)
	if err != nil && l.Errors != nil {
		l.Errors <- err
	}
}

// --------------------------------------------------------------------------------
// finalizers
// --------------------------------------------------------------------------------

// Close releases shared resources for the agent.
// It will stop listeners and wait for them to complete work
// and then zero out any other resources.
func (l *Logger) Close() {
	l.Lock()
	defer l.Unlock()

	if l.Flags != nil {
		l.Flags.SetNone()
	}

	for _, listeners := range l.Listeners {
		for _, listener := range listeners {
			_ = listener.Stop()
		}
	}
	if closer, ok := l.Output.(io.Closer); ok {
		_ = closer.Close()
	}
	l.Listeners = nil
	l.Filters = nil
}

// Drain stops the event listeners, letting them complete their work
// and then restarts the listeners.
func (l *Logger) Drain() {
	l.DrainContext(context.Background())
}

// DrainContext waits for the logger to finish its queue of events with a given context.
func (l *Logger) DrainContext(ctx context.Context) {
	for _, workers := range l.Listeners {
		for _, worker := range workers {
			_ = worker.StopContext(ctx)
			worker.Reset()

			notifyStarted := worker.NotifyStarted()
			go func(w *Worker) {
				_ = w.Start()
			}(worker)

			// Wait for worker to start
			select {
			case <-notifyStarted:
			case <-ctx.Done():
			}
		}
	}
}

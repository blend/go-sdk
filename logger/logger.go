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
		Formatter:     NewTextOutputFormatter(),
		Output:        NewInterlockedWriter(os.Stdout),
		RecoverPanics: DefaultRecoverPanics,
		Flags:         NewFlags(DefaultFlags...),
	}
	l.Context = NewContext(l, nil, nil)
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
	return MustNew(append([]Option{
		OptConfigFromEnv(),
		OptAll(),
	}, options...)...)
}

// None returns a new logger with all flags enabled.
func None() *Logger {
	return MustNew(OptNone(), OptOutput(nil))
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

// Logger is a handler for various logging events with descendent handlers.
type Logger struct {
	sync.Mutex
	*Flags
	Context

	RecoverPanics bool

	Output    io.Writer
	Formatter WriteFormatter
	Errors    chan error
	Listeners map[string]map[string]*Worker
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

// HasListener returns if a specific listener is registerd for a flag.
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

	w := NewWorker(listener)
	if listeners, ok := l.Listeners[flag]; ok {
		listeners[listenerName] = w
	} else {
		l.Listeners[flag] = map[string]*Worker{
			listenerName: w,
		}
	}
	go w.Start()
	<-w.NotifyStarted()
}

// RemoveListeners clears *all* listeners for a Flag.
func (l *Logger) RemoveListeners(flag string) {
	l.Lock()
	defer l.Unlock()

	if l.Listeners == nil {
		return
	}

	listeners, ok := l.Listeners[flag]
	if !ok {
		return
	}

	for _, l := range listeners {
		l.Stop()
	}

	delete(l.Listeners, flag)
}

// RemoveListener clears a specific listener for a Flag.
func (l *Logger) RemoveListener(flag, listenerName string) {
	l.Lock()
	defer l.Unlock()

	if l.Listeners == nil {
		return
	}

	listeners, ok := l.Listeners[flag]
	if !ok {
		return
	}

	worker, ok := listeners[listenerName]
	if !ok {
		return
	}

	worker.Stop()
	<-worker.NotifyStopped()

	delete(listeners, listenerName)
	if len(listeners) == 0 {
		delete(l.Listeners, flag)
	}
}

// Trigger fires the listeners for a given event asynchronously, and writes the event to the output.
// The invocations will be queued in a work queue per listener.
// There are no order guarantees on when these events will be processed across listeners.
// This call will not block on the event listeners, but will block on the write.
func (l *Logger) trigger(ctx context.Context, e Event, sync bool) {
	if e == nil {
		return
	}

	flag := e.GetFlag()
	if !l.IsEnabled(flag) {
		return
	}

	if !IsSkipTrigger(ctx) {
		var listeners map[string]*Worker
		l.Lock()
		if l.Listeners != nil {
			if flagListeners, ok := l.Listeners[flag]; ok {
				listeners = flagListeners
			}
		}
		l.Unlock()

		for _, listener := range listeners {
			if sync {
				listener.Process(EventWithContext{ctx, e})
			} else {
				listener.Work <- EventWithContext{ctx, e}
			}
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

	err := l.Formatter.WriteFormat(ctx, l.Output, e)
	if err != nil && l.Errors != nil {
		l.Errors <- err
	}
}

// --------------------------------------------------------------------------------
// finalizers
// --------------------------------------------------------------------------------

// Close releases shared resources for the agent.
func (l *Logger) Close() error {
	l.Lock()
	defer l.Unlock()

	if l.Flags != nil {
		l.Flags.SetNone()
	}

	for _, listeners := range l.Listeners {
		for _, listener := range listeners {
			listener.Stop()
		}
	}

	for key := range l.Listeners {
		delete(l.Listeners, key)
	}
	l.Listeners = nil
	return nil
}

// Drain waits for the logger to finish its queue of events.
func (l *Logger) Drain() error {
	return l.DrainContext(context.Background())
}

// DrainContext waits for the logger to finish its queue of events with a given context.
func (l *Logger) DrainContext(ctx context.Context) error {
	var err error
	for _, workers := range l.Listeners {
		for _, worker := range workers {
			if err = worker.DrainContext(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

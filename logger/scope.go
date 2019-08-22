package logger

import (
	"context"
	"fmt"
	"time"
)

var (
	_ Log = (*Scope)(nil)
)

// NewScope returns a new scope for a logger with a given set of optional options.
func NewScope(log *Logger, options ...ScopeOption) Scope {
	s := Scope{
		Logger:      log,
		Context:     context.Background(),
		Labels:      Labels{},
		Annotations: Annotations{},
	}
	for _, option := range options {
		option(&s)
	}
	return s
}

// Scope is a logger scope.
// It is used to split a logger into functional concerns but retain all the underlying functionality of logging.
// You can attach extra data (Fields) to the scope (useful for things like the Environment).
// You can also set a context to be used when triggering events.
type Scope struct {
	// Path is a series of descriptive labels that shows the origin of the scope.
	Path []string

	// Labels are descriptive string fields for the scope.
	Labels
	// Annotations are extra fields for the scope.
	Annotations

	// Context is a relevant context for the scope, it is passed to listeners for events.
	// Before triggering events, it is loaded with the Path and Fields from the Scope as Values.
	Context context.Context
	// Logger is a parent reference to the root logger; this holds
	// information around what flags are enabled and listeners for events.
	Logger *Logger
}

// ScopeOption is a mutator for a scope.
type ScopeOption func(*Scope)

// OptScopePath sets the path on a scope.
func OptScopePath(path ...string) ScopeOption {
	return func(s *Scope) {
		s.Path = path
	}
}

// OptScopeLabels sets the labels on a scope.
func OptScopeLabels(labels ...Labels) ScopeOption {
	return func(s *Scope) {
		s.Labels = CombineLabels(labels...)
	}
}

// OptScopeAnnotations sets the annotations on a scope.
func OptScopeAnnotations(annotations ...Annotations) ScopeOption {
	return func(s *Scope) {
		s.Annotations = CombineAnnotations(annotations...)
	}
}

// OptScopeContext sets the context on a scope.
// This context will be used as the triggering context for any events.
func OptScopeContext(ctx context.Context) ScopeOption {
	return func(s *Scope) {
		s.Context = ctx
	}
}

// WithContext returns a new scope context.
func (sc Scope) WithContext(ctx context.Context) Scope {
	return NewScope(sc.Logger,
		OptScopePath(sc.Path...),
		OptScopeLabels(sc.Labels),
		OptScopeAnnotations(sc.Annotations),
		OptScopeContext(ctx),
	)
}

// WithPath returns a new scope with a given additional path segment.
func (sc Scope) WithPath(paths ...string) Scope {
	return NewScope(sc.Logger,
		OptScopePath(append(sc.Path, paths...)...),
		OptScopeLabels(sc.Labels),
		OptScopeAnnotations(sc.Annotations),
		OptScopeContext(sc.Context),
	)
}

// WithLabels returns a new scope with a given additional set of labels.
func (sc Scope) WithLabels(labels Labels) Scope {
	return NewScope(sc.Logger,
		OptScopePath(sc.Path...),
		OptScopeLabels(sc.Labels, labels),
		OptScopeAnnotations(sc.Annotations),
		OptScopeContext(sc.Context),
	)
}

// WithAnnotations returns a new scope with a given additional set of annotations.
func (sc Scope) WithAnnotations(annotations Annotations) Scope {
	return NewScope(sc.Logger,
		OptScopePath(sc.Path...),
		OptScopeLabels(sc.Labels),
		OptScopeAnnotations(sc.Annotations, annotations),
		OptScopeContext(sc.Context),
	)
}

// --------------------------------------------------------------------------------
// Trigger event handler
// --------------------------------------------------------------------------------

// Trigger triggers an event in the subcontext.
func (sc Scope) Trigger(ctx context.Context, event Event) {
	sc.Logger.Trigger(sc.ApplyContext(ctx), event)
}

// --------------------------------------------------------------------------------
// Builtin Flag Handlers (infof, debugf etc.)
// --------------------------------------------------------------------------------

// Info logs an informational message to the output stream.
func (sc Scope) Info(args ...interface{}) {
	sc.Trigger(sc.Context, NewMessageEvent(Info, fmt.Sprint(args...)))
}

// Infof logs an informational message to the output stream.
func (sc Scope) Infof(format string, args ...interface{}) {
	sc.Trigger(sc.Context, NewMessageEvent(Info, fmt.Sprintf(format, args...)))
}

// Debug logs a debug message to the output stream.
func (sc Scope) Debug(args ...interface{}) {
	sc.Trigger(sc.Context, NewMessageEvent(Debug, fmt.Sprint(args...)))
}

// Debugf logs a debug message to the output stream.
func (sc Scope) Debugf(format string, args ...interface{}) {
	sc.Trigger(sc.Context, NewMessageEvent(Debug, fmt.Sprintf(format, args...)))
}

// Warningf logs a warning message to the output stream.
func (sc Scope) Warningf(format string, args ...interface{}) {
	sc.Trigger(sc.Context, NewErrorEvent(Warning, fmt.Errorf(format, args...)))
}

// Errorf writes an event to the log and triggers event listeners.
func (sc Scope) Errorf(format string, args ...interface{}) {
	sc.Trigger(sc.Context, NewErrorEvent(Error, fmt.Errorf(format, args...)))
}

// Fatalf writes an event to the log and triggers event listeners.
func (sc Scope) Fatalf(format string, args ...interface{}) {
	sc.Trigger(sc.Context, NewErrorEvent(Fatal, fmt.Errorf(format, args...)))
}

// Warning logs a warning error to std err.
func (sc Scope) Warning(err error, opts ...ErrorEventOption) error {
	sc.Trigger(sc.Context, NewErrorEvent(Warning, err, opts...))
	return err
}

// Error logs an error to std err.
func (sc Scope) Error(err error, opts ...ErrorEventOption) error {
	sc.Trigger(sc.Context, NewErrorEvent(Error, err, opts...))
	return err
}

// Fatal logs an error as fatal.
func (sc Scope) Fatal(err error, opts ...ErrorEventOption) error {
	sc.Trigger(sc.Context, NewErrorEvent(Fatal, err, opts...))
	return err
}

// ApplyContext applies the scope context to a given context.
func (sc Scope) ApplyContext(ctx context.Context) context.Context {
	ctx = WithTriggerTimestamp(ctx, time.Now().UTC())
	ctx = WithScopePath(ctx, append(sc.Path, GetScopePath(ctx)...)...)
	ctx = WithLabels(ctx, CombineLabels(sc.Labels, GetLabels(ctx)))
	ctx = WithAnnotations(ctx, CombineAnnotations(sc.Annotations, GetAnnotations(ctx)))
	return ctx
}

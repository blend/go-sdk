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
		Labels:      make(Labels),
		Annotations: make(Annotations),
	}
	for _, option := range options {
		option(&s)
	}
	return s
}

// Scope is a set of re-usable parameters for triggering events.
/*
The key fields:
- "Path" is a set of names that denote a hierarchy or tree of calls.
- "Labels" are string pairs that will appear with written log messages for easier searching later.
- "Annoations" are string pairs taht will not appear with written log messages, but can be used to add extra data to events.
*/
type Scope struct {
	// Path is a series of descriptive labels that shows the origin of the scope.
	Path []string
	// Labels are descriptive string fields for the scope.
	Labels
	// Annotations are extra fields for the scope.
	Annotations

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

// WithPath returns a new scope with a given additional path segment.
func (sc Scope) WithPath(paths ...string) Scope {
	return NewScope(sc.Logger,
		OptScopePath(append(sc.Path, paths...)...),
		OptScopeLabels(sc.Labels),
		OptScopeAnnotations(sc.Annotations),
	)
}

// WithLabels returns a new scope with a given additional set of labels.
func (sc Scope) WithLabels(labels Labels) Scope {
	return NewScope(sc.Logger,
		OptScopePath(sc.Path...),
		OptScopeLabels(sc.Labels, labels),
		OptScopeAnnotations(sc.Annotations),
	)
}

// WithAnnotations returns a new scope with a given additional set of annotations.
func (sc Scope) WithAnnotations(annotations Annotations) Scope {
	return NewScope(sc.Logger,
		OptScopePath(sc.Path...),
		OptScopeLabels(sc.Labels),
		OptScopeAnnotations(sc.Annotations, annotations),
	)
}

// --------------------------------------------------------------------------------
// Trigger event handler
// --------------------------------------------------------------------------------

// Trigger triggers an event in the subcontext.
// The provided context is ammended with fields from the scope.
// The provided context is also ammended with a TriggerTimestamp, which can be retrieved with `GetTriggerTimestamp(ctx)` in listeners.
func (sc Scope) Trigger(event Event) {
	sc.TriggerContext(context.Background(), event)
}

// TriggerContext triggers an event with a given context..
func (sc Scope) TriggerContext(ctx context.Context, event Event) {
	ctx = WithTriggerTimestamp(ctx, time.Now().UTC())
	sc.Logger.Dispatch(sc.ApplyContext(ctx), event)
}

// --------------------------------------------------------------------------------
// Builtin Flag Handlers (infof, debugf etc.)
// --------------------------------------------------------------------------------

// Info logs an informational message to the output stream.
func (sc Scope) Info(args ...interface{}) {
	sc.Trigger(NewMessageEvent(Info, fmt.Sprint(args...)))
}

// InfoContext logs an informational message to the output stream in a given context.
func (sc Scope) InfoContext(ctx context.Context, args ...interface{}) {
	sc.TriggerContext(ctx, NewMessageEvent(Info, fmt.Sprint(args...)))
}

// Infof logs an informational message to the output stream.
func (sc Scope) Infof(format string, args ...interface{}) {
	sc.Trigger(NewMessageEvent(Info, fmt.Sprintf(format, args...)))
}

// InfofContext logs an informational message to the output stream in a given context.
func (sc Scope) InfofContext(ctx context.Context, format string, args ...interface{}) {
	sc.TriggerContext(ctx, NewMessageEvent(Info, fmt.Sprintf(format, args...)))
}

// Debug logs a debug message to the output stream.
func (sc Scope) Debug(args ...interface{}) {
	sc.Trigger(NewMessageEvent(Debug, fmt.Sprint(args...)))
}

// DebugContext logs a debug message to the output stream in a given context.
func (sc Scope) DebugContext(ctx context.Context, args ...interface{}) {
	sc.TriggerContext(ctx, NewMessageEvent(Debug, fmt.Sprint(args...)))
}

// Debugf logs a debug message to the output stream.
func (sc Scope) Debugf(format string, args ...interface{}) {
	sc.Trigger(NewMessageEvent(Debug, fmt.Sprintf(format, args...)))
}

// DebugfContext logs a debug message to the output stream.
func (sc Scope) DebugfContext(ctx context.Context, format string, args ...interface{}) {
	sc.TriggerContext(ctx, NewMessageEvent(Debug, fmt.Sprintf(format, args...)))
}

// Warningf logs a warning message to the output stream.
func (sc Scope) Warningf(format string, args ...interface{}) {
	sc.Trigger(NewErrorEvent(Warning, fmt.Errorf(format, args...)))
}

// WarningfContext logs a warning message to the output stream in a given context.
func (sc Scope) WarningfContext(ctx context.Context, format string, args ...interface{}) {
	sc.TriggerContext(ctx, NewErrorEvent(Warning, fmt.Errorf(format, args...)))
}

// Errorf writes an event to the log and triggers event listeners.
func (sc Scope) Errorf(format string, args ...interface{}) {
	sc.Trigger(NewErrorEvent(Error, fmt.Errorf(format, args...)))
}

// ErrorfContext writes an event to the log and triggers event listeners in a given context.
func (sc Scope) ErrorfContext(ctx context.Context, format string, args ...interface{}) {
	sc.TriggerContext(ctx, NewErrorEvent(Error, fmt.Errorf(format, args...)))
}

// Fatalf writes an event to the log and triggers event listeners.
func (sc Scope) Fatalf(format string, args ...interface{}) {
	sc.Trigger(NewErrorEvent(Fatal, fmt.Errorf(format, args...)))
}

// FatalfContext writes an event to the log and triggers event listeners in a given context.
func (sc Scope) FatalfContext(ctx context.Context, format string, args ...interface{}) {
	sc.TriggerContext(ctx, NewErrorEvent(Fatal, fmt.Errorf(format, args...)))
}

// Warning logs a warning error to std err.
func (sc Scope) Warning(err error, opts ...ErrorEventOption) {
	sc.Trigger(NewErrorEvent(Warning, err, opts...))
}

// WarningContext logs a warning error to std err in a given context.
func (sc Scope) WarningContext(ctx context.Context, err error, opts ...ErrorEventOption) {
	sc.TriggerContext(ctx, NewErrorEvent(Warning, err, opts...))
}

// Error logs an error to std err.
func (sc Scope) Error(err error, opts ...ErrorEventOption) {
	sc.Trigger(NewErrorEvent(Error, err, opts...))
}

// ErrorContext logs an error to std err.
func (sc Scope) ErrorContext(ctx context.Context, err error, opts ...ErrorEventOption) {
	sc.TriggerContext(ctx, NewErrorEvent(Error, err, opts...))
}

// Fatal logs an error as fatal.
func (sc Scope) Fatal(err error, opts ...ErrorEventOption) {
	sc.Trigger(NewErrorEvent(Fatal, err, opts...))
}

// FatalContext logs an error as fatal.
func (sc Scope) FatalContext(ctx context.Context, err error, opts ...ErrorEventOption) {
	sc.TriggerContext(ctx, NewErrorEvent(Fatal, err, opts...))
}

//
// Context utilities
//

// FromContext returns a scope from a given context.
// It will read any relevant fields off the context (Path, Labels, Annotations)
// and append them to values already on the scope.
func (sc Scope) FromContext(ctx context.Context) Scope {
	return NewScope(sc.Logger,
		OptScopePath(append(sc.Path, GetPath(ctx)...)...),
		OptScopeLabels(sc.Labels, GetLabels(ctx)),
		OptScopeAnnotations(sc.Annotations, GetAnnotations(ctx)),
	)
}

// ApplyContext applies the scope fields to a given context.
func (sc Scope) ApplyContext(ctx context.Context) context.Context {
	ctx = WithPath(ctx, append(sc.Path, GetPath(ctx)...)...)
	ctx = WithLabels(ctx, CombineLabels(sc.Labels, GetLabels(ctx)))
	ctx = WithAnnotations(ctx, CombineAnnotations(sc.Annotations, GetAnnotations(ctx)))
	return ctx
}

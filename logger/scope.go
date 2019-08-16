package logger

import (
	"context"
	"fmt"
	"net/http"
)

var (
	_ Log = (*Scope)(nil)
)

// NewScope returns a new scope for a logger with a given set of optional options.
func NewScope(log *Logger, options ...ScopeOption) Scope {
	s := Scope{
		Logger:  log,
		Context: context.Background(),
		Fields:  Fields{},
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
	// Fields are descriptive fields for the scope.
	Fields Fields
	// Context is a relevant context for the scope, it is passed to listeners for events.
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

// OptScopeFields sets the fields on a scope.
func OptScopeFields(fields ...Fields) ScopeOption {
	return func(s *Scope) {
		s.Fields = CombineFields(fields...)
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
		OptScopeFields(sc.Fields),
		OptScopeContext(ctx),
	)
}

// WithPath returns a new scope with a given additional path segment.
func (sc Scope) WithPath(paths ...string) Scope {
	return NewScope(sc.Logger,
		OptScopeContext(sc.Context),
		OptScopePath(append(sc.Path, paths...)...),
		OptScopeFields(sc.Fields),
	)
}

// WithFields returns a new scope with a given additional set of fields.
func (sc Scope) WithFields(fields Fields) Scope {
	return NewScope(sc.Logger,
		OptScopePath(sc.Path...),
		OptScopeFields(sc.Fields, fields),
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
func (sc Scope) Warning(err error) error {
	sc.Trigger(sc.Context, NewErrorEvent(Warning, err))
	return err
}

// WarningWithReq logs a warning error to std err with a request.
func (sc Scope) WarningWithReq(err error, req *http.Request) error {
	ee := NewErrorEvent(Warning, err)
	ee.Request = req
	sc.Trigger(sc.Context, ee)
	return err
}

// Error logs an error to std err.
func (sc Scope) Error(err error) error {
	sc.Trigger(sc.Context, NewErrorEvent(Error, err))
	return err
}

// ErrorWithReq logs an error to std err with a request.
func (sc Scope) ErrorWithReq(err error, req *http.Request) error {
	ee := NewErrorEvent(Error, err)
	ee.Request = req
	sc.Trigger(sc.Context, ee)
	return err
}

// Fatal logs an error as fatal.
func (sc Scope) Fatal(err error) error {
	sc.Trigger(sc.Context, NewErrorEvent(Fatal, err))
	return err
}

// FatalWithReq logs an error as fatal with a request as state.
func (sc Scope) FatalWithReq(err error, req *http.Request) error {
	ee := NewErrorEvent(Fatal, err)
	ee.Request = req
	sc.Trigger(sc.Context, ee)
	return err
}

// ApplyContext applies the scope context to a given context.
func (sc Scope) ApplyContext(ctx context.Context) context.Context {
	ctx = WithScopePath(ctx, sc.Path...)
	ctx = WithFields(ctx, CombineFields(sc.Fields, GetFields(ctx)))
	return ctx
}

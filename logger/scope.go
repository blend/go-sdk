package logger

import (
	"context"
	"fmt"
	"net/http"
)

var (
	_ Scoper = (*Scope)(nil)
)

// NewScope returns a new context.
func NewScope(ctx context.Context, log *Logger) Scope {
	return Scope{
		ctx:    ctx,
		Logger: log,
	}
}

// Scope is a logger context.
// It is used to split a logger into functional concerns but retain all the underlying functionality of logging.
type Scope struct {
	// ctx is the metadata we should trigger the event with.
	ctx context.Context
	// Logger is a parent reference to the root logger.
	Logger *Logger
}

// Context returns the underlying context.
func (sc Scope) Context() context.Context {
	return sc.ctx
}

// WithContext returns a new scope context.
func (sc Scope) WithContext(ctx context.Context) Scope {
	return NewScope(sc.ApplyContext(ctx), sc.Logger)
}

// WithPath returns a new scope with a given additional path segment.
func (sc Scope) WithPath(paths ...string) Scope {
	return NewScope(WithScopePath(sc.ctx, paths...), sc.Logger)
}

// WithLabels returns a new scope with a given additional set of labels.
func (sc Scope) WithLabels(values Labels) Scope {
	return NewScope(WithLabels(sc.ctx, values), sc.Logger)
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
	sc.Logger.Trigger(sc.ctx, NewMessageEvent(Info, fmt.Sprint(args...)))
}

// Infof logs an informational message to the output stream.
func (sc Scope) Infof(format string, args ...interface{}) {
	sc.Logger.Trigger(sc.ctx, NewMessageEvent(Info, fmt.Sprintf(format, args...)))
}

// Debug logs a debug message to the output stream.
func (sc Scope) Debug(args ...interface{}) {
	sc.Logger.Trigger(sc.ctx, NewMessageEvent(Debug, fmt.Sprint(args...)))
}

// Debugf logs a debug message to the output stream.
func (sc Scope) Debugf(format string, args ...interface{}) {
	sc.Logger.Trigger(sc.ctx, NewMessageEvent(Debug, fmt.Sprintf(format, args...)))
}

// Warningf logs a warning message to the output stream.
func (sc Scope) Warningf(format string, args ...interface{}) {
	sc.Logger.Trigger(sc.ctx, NewErrorEvent(Warning, fmt.Errorf(format, args...)))
}

// Errorf writes an event to the log and triggers event listeners.
func (sc Scope) Errorf(format string, args ...interface{}) {
	sc.Logger.Trigger(sc.ctx, NewErrorEvent(Error, fmt.Errorf(format, args...)))
}

// Fatalf writes an event to the log and triggers event listeners.
func (sc Scope) Fatalf(format string, args ...interface{}) {
	sc.Logger.Trigger(sc.ctx, NewErrorEvent(Fatal, fmt.Errorf(format, args...)))
}

// Warning logs a warning error to std err.
func (sc Scope) Warning(err error) error {
	sc.Logger.Trigger(sc.ctx, NewErrorEvent(Warning, err))
	return err
}

// WarningWithReq logs a warning error to std err with a request.
func (sc Scope) WarningWithReq(err error, req *http.Request) error {
	ee := NewErrorEvent(Warning, err)
	ee.Request = req
	sc.Logger.Trigger(sc.ctx, ee)
	return err
}

// Error logs an error to std err.
func (sc Scope) Error(err error) error {
	sc.Logger.Trigger(sc.ctx, NewErrorEvent(Error, err))
	return err
}

// ErrorWithReq logs an error to std err with a request.
func (sc Scope) ErrorWithReq(err error, req *http.Request) error {
	ee := NewErrorEvent(Error, err)
	ee.Request = req
	sc.Logger.Trigger(sc.ctx, ee)
	return err
}

// Fatal logs an error as fatal.
func (sc Scope) Fatal(err error) error {
	sc.Logger.Trigger(sc.ctx, NewErrorEvent(Fatal, err))
	return err
}

// FatalWithReq logs an error as fatal with a request as state.
func (sc Scope) FatalWithReq(err error, req *http.Request) error {
	ee := NewErrorEvent(Fatal, err)
	ee.Request = req
	sc.Logger.Trigger(sc.ctx, ee)
	return err
}

// ApplyContext applies the scope context to a given context.
func (sc Scope) ApplyContext(ctx context.Context) context.Context {
	ctx = WithScopePath(ctx, GetScopePath(sc.ctx)...)
	ctx = WithLabels(ctx, GetLabels(sc.ctx))
	return ctx
}

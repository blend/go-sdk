package logger

import (
	"context"
	"fmt"
	"net/http"
)

// NewScope returns a new context.
func NewScope(log *Logger, path []string, labels Labels, annotations Annotations, opts ...ScopeOption) Scope {
	c := Scope{
		Logger:      log,
		Path:        path,
		Labels:      labels,
		Annotations: annotations,
	}
	for _, opt := range opts {
		opt(&c)
	}
	return c
}

// ScopeOption is an option for contexts.
type ScopeOption func(*Scope)

// OptScopePath appends new path segments to the context.
func OptScopePath(path ...string) ScopeOption {
	return func(c *Scope) {
		c.Path = append(c.Path, path...)
	}
}

// OptScopeSetPath sets path on the context.
func OptScopeSetPath(path ...string) ScopeOption {
	return func(c *Scope) {
		c.Path = path
	}
}

// OptScopeLabels adds fields to the context.
func OptScopeLabels(labels Labels) ScopeOption {
	return func(c *Scope) {
		c.Labels = CombineLabels(c.Labels, labels)
	}
}

// OptScopeSetLabels sets fields on the context.
func OptScopeSetLabels(labels Labels) ScopeOption {
	return func(c *Scope) {
		c.Labels = labels
	}
}

// Scope is a logger context.
// It is used to split a logger into functional concerns
// but retain all the underlying machinery of logging.
type Scope struct {
	Labels
	Annotations

	Logger *Logger
	Path   []string
}

// SubScope returns a new sub context.
func (sc Scope) SubScope(name string, options ...ScopeOption) Scope {
	return NewScope(sc.Logger, append(sc.Path, name), sc.Labels, sc.Annotations, options...)
}

// WithLabels returns a new sub context.
func (sc Scope) WithLabels(labels Labels, options ...ScopeOption) Scope {
	return NewScope(sc.Logger, sc.Path, CombineLabels(sc.Labels, labels), sc.Annotations, options...)
}

// WithAnnotations returns a new sub context.
func (sc Scope) WithAnnotations(annotations Annotations, options ...ScopeOption) Scope {
	return NewScope(sc.Logger, sc.Path, sc.Labels, CombineAnnotations(sc.Annotations, annotations), options...)
}

// --------------------------------------------------------------------------------
// Trigger event handler
// --------------------------------------------------------------------------------

// Trigger triggers an event in the subcontext.
func (sc Scope) Trigger(ctx context.Context, event Event) {
	sc.Logger.trigger(WithSubScopeMeta(ctx, sc), event, false)
}

// SyncTrigger triggers an event in the subcontext synchronously..
func (sc Scope) SyncTrigger(ctx context.Context, event Event) {
	sc.Logger.trigger(WithSubScopeMeta(ctx, sc), event, true)
}

// --------------------------------------------------------------------------------
// Builtin Flag Handlers (infof, debugf etc.)
// --------------------------------------------------------------------------------

// Info logs an informational message to the output stream.
func (sc Scope) Info(args ...interface{}) {
	sc.Trigger(context.Background(), NewMessageEvent(Info, fmt.Sprint(args...)))
}

// Infof logs an informational message to the output stream.
func (sc Scope) Infof(format string, args ...interface{}) {
	sc.Trigger(context.Background(), NewMessageEvent(Info, fmt.Sprintf(format, args...)))
}

// Debug logs a debug message to the output stream.
func (sc Scope) Debug(args ...interface{}) {
	sc.Trigger(context.Background(), NewMessageEvent(Debug, fmt.Sprint(args...)))
}

// Debugf logs a debug message to the output stream.
func (sc Scope) Debugf(format string, args ...interface{}) {
	sc.Trigger(context.Background(), NewMessageEvent(Debug, fmt.Sprintf(format, args...)))
}

// Warningf logs a warning message to the output stream.
func (sc Scope) Warningf(format string, args ...interface{}) {
	sc.Trigger(context.Background(), NewErrorEvent(Warning, fmt.Errorf(format, args...)))
}

// Errorf writes an event to the log and triggers event listeners.
func (sc Scope) Errorf(format string, args ...interface{}) {
	sc.Trigger(context.Background(), NewErrorEvent(Error, fmt.Errorf(format, args...)))
}

// Fatalf writes an event to the log and triggers event listeners.
func (sc Scope) Fatalf(format string, args ...interface{}) {
	sc.Trigger(context.Background(), NewErrorEvent(Fatal, fmt.Errorf(format, args...)))
}

// Warning logs a warning error to std err.
func (sc Scope) Warning(err error) error {
	sc.Trigger(context.Background(), NewErrorEvent(Warning, err))
	return err
}

// WarningWithReq logs a warning error to std err with a request.
func (sc Scope) WarningWithReq(err error, req *http.Request) error {
	ee := NewErrorEvent(Warning, err)
	ee.State = req
	sc.Trigger(context.Background(), ee)
	return err
}

// Error logs an error to std err.
func (sc Scope) Error(err error) error {
	sc.Trigger(context.Background(), NewErrorEvent(Error, err))
	return err
}

// ErrorWithReq logs an error to std err with a request.
func (sc Scope) ErrorWithReq(err error, req *http.Request) error {
	ee := NewErrorEvent(Error, err)
	ee.State = req
	sc.Trigger(context.Background(), ee)
	return err
}

// Fatal logs an error as fatal.
func (sc Scope) Fatal(err error) error {
	sc.Trigger(context.Background(), NewErrorEvent(Fatal, err))
	return err
}

// FatalWithReq logs an error as fatal with a request as state.
func (sc Scope) FatalWithReq(err error, req *http.Request) error {
	ee := NewErrorEvent(Fatal, err)
	ee.State = req
	sc.Trigger(context.Background(), ee)
	return err
}

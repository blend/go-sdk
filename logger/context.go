package logger

import (
	"context"
	"fmt"
	"net/http"
)

// NewContext returns a new context.
func NewContext(log *Logger, path []string, fields Fields, opts ...ContextOption) Context {
	c := Context{
		Logger: log,
		Path:   path,
		Fields: fields,
	}
	for _, opt := range opts {
		opt(&c)
	}
	return c
}

// ContextOption is an option for contexts.
type ContextOption func(*Context)

// OptContextPath appends new path segments to the context.
func OptContextPath(path ...string) ContextOption {
	return func(c *Context) {
		c.Path = append(c.Path, path...)
	}
}

// OptContextSetPath sets path on the context.
func OptContextSetPath(path ...string) ContextOption {
	return func(c *Context) {
		c.Path = path
	}
}

// OptContextFields adds fields to the context.
func OptContextFields(fields Fields) ContextOption {
	return func(c *Context) {
		c.Fields = CombineFields(c.Fields, fields)
	}
}

// OptContextSetFields sets fields on the context.
func OptContextSetFields(fields Fields) ContextOption {
	return func(c *Context) {
		c.Fields = fields
	}
}

// Context is a logger context.
// It is used to split a logger into functional concerns
// but retain all the underlying machinery of logging.
type Context struct {
	Logger *Logger
	Path   []string
	Fields Fields
}

// SubContext returns a new sub context.
func (sc Context) SubContext(name string, options ...ContextOption) Context {
	return NewContext(sc.Logger, append(sc.Path, name), sc.Fields, options...)
}

// WithFields returns a new sub context.
func (sc Context) WithFields(fields Fields, options ...ContextOption) Context {
	return NewContext(sc.Logger, sc.Path, CombineFields(sc.Fields, fields), options...)
}

// --------------------------------------------------------------------------------
// Trigger event handler
// --------------------------------------------------------------------------------

// Trigger triggers an event in the subcontext.
func (sc Context) Trigger(ctx context.Context, event Event) {
	sc.Logger.trigger(WithSubContextMeta(ctx, sc.Path, sc.Fields), event, false)
}

// SyncTrigger triggers an event in the subcontext synchronously..
func (sc Context) SyncTrigger(ctx context.Context, event Event) {
	sc.Logger.trigger(WithSubContextMeta(ctx, sc.Path, sc.Fields), event, true)
}

// --------------------------------------------------------------------------------
// Builtin Flag Handlers (infof, debugf etc.)
// --------------------------------------------------------------------------------

// Info logs an informational message to the output stream.
func (sc Context) Info(args ...interface{}) {
	sc.Trigger(context.Background(), NewMessageEvent(Info, fmt.Sprint(args...)))
}

// Infof logs an informational message to the output stream.
func (sc Context) Infof(format string, args ...interface{}) {
	sc.Trigger(context.Background(), NewMessageEvent(Info, fmt.Sprintf(format, args...)))
}

// Debug logs a debug message to the output stream.
func (sc Context) Debug(args ...interface{}) {
	sc.Trigger(context.Background(), NewMessageEvent(Debug, fmt.Sprint(args...)))
}

// Debugf logs a debug message to the output stream.
func (sc Context) Debugf(format string, args ...interface{}) {
	sc.Trigger(context.Background(), NewMessageEvent(Debug, fmt.Sprintf(format, args...)))
}

// Warningf logs a warning message to the output stream.
func (sc Context) Warningf(format string, args ...interface{}) {
	sc.Trigger(context.Background(), NewErrorEvent(Warning, fmt.Errorf(format, args...)))
}

// Errorf writes an event to the log and triggers event listeners.
func (sc Context) Errorf(format string, args ...interface{}) {
	sc.Trigger(context.Background(), NewErrorEvent(Error, fmt.Errorf(format, args...)))
}

// Fatalf writes an event to the log and triggers event listeners.
func (sc Context) Fatalf(format string, args ...interface{}) {
	sc.Trigger(context.Background(), NewErrorEvent(Fatal, fmt.Errorf(format, args...)))
}

// Warning logs a warning error to std err.
func (sc Context) Warning(err error) error {
	sc.Trigger(context.Background(), NewErrorEvent(Warning, err))
	return err
}

// WarningWithReq logs a warning error to std err with a request.
func (sc Context) WarningWithReq(err error, req *http.Request) error {
	ee := NewErrorEvent(Warning, err)
	ee.State = req
	sc.Trigger(context.Background(), ee)
	return err
}

// Error logs an error to std err.
func (sc Context) Error(err error) error {
	sc.Trigger(context.Background(), NewErrorEvent(Error, err))
	return err
}

// ErrorWithReq logs an error to std err with a request.
func (sc Context) ErrorWithReq(err error, req *http.Request) error {
	ee := NewErrorEvent(Error, err)
	ee.State = req
	sc.Trigger(context.Background(), ee)
	return err
}

// Fatal logs an error as fatal.
func (sc Context) Fatal(err error) error {
	sc.Trigger(context.Background(), NewErrorEvent(Fatal, err))
	return err
}

// FatalWithReq logs an error as fatal with a request as state.
func (sc Context) FatalWithReq(err error, req *http.Request) error {
	ee := NewErrorEvent(Fatal, err)
	ee.State = req
	sc.Trigger(context.Background(), ee)
	return err
}

/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import "context"

// Listenable is an interface loggers can ascribe to.
type Listenable interface {
	Listen(flag string, label string, listener Listener)
}

// Filterable is an interface loggers can ascribe to.
type Filterable interface {
	Filter(flag string, label string, filter Filter)
}

// FilterListenable is a type that loggers can ascribe to.
type FilterListenable interface {
	Filterable
	Listenable
}

// Triggerable is type that can trigger events.
type Triggerable interface {
	// TriggerContext fires listeners for an event with a given context.
	TriggerContext(context.Context, Event)
}

// Scoper is a type that can return a scope.
type Scoper interface {
	// FromContext returns a new scope, merging the existing scope fields with fields found
	// on a given context.
	FromContext(context.Context) Scope
	// Apply augments a given context with fields from the Scope, including Labels, Annotations, and Path.
	ApplyContext(context.Context) context.Context
	// WithPath returns a new scope with a given set of additional path segments.
	WithPath(...string) Scope
	// WithLabels returns a new scope with a given set of additional label values.
	WithLabels(Labels) Scope
	// WithAnnotations returns a new scope with a given set of additional annotation values.
	WithAnnotations(Annotations) Scope
}

// Writable is a type that can write events.
type Writable interface {
	Write(context.Context, Event)
}

// ApplyContexter is a type that modifies a context.
type ApplyContexter interface {
	ApplyContext(context.Context) context.Context
}

// WriteTriggerable is a type that can both trigger and write events.
type WriteTriggerable interface {
	Writable
	Triggerable
}

// InfoReceiver is a type that defines Info.
type InfoReceiver interface {
	Info(...interface{})
	InfoContext(context.Context, ...interface{})
}

// PrintfReceiver is a type that defines Printf.
type PrintfReceiver interface {
	Printf(string, ...interface{})
}

// PrintReceiver is a type that defines Print.
type PrintReceiver interface {
	Print(...interface{})
}

// PrintlnReceiver is a type that defines Println.
type PrintlnReceiver interface {
	Println(...interface{})
}

// InfofReceiver is a type that defines Infof.
type InfofReceiver interface {
	Infof(string, ...interface{})
	InfofContext(context.Context, string, ...interface{})
}

// DebugReceiver is a type that defines Debug.
type DebugReceiver interface {
	Debug(...interface{})
	DebugContext(context.Context, ...interface{})
}

// DebugfReceiver is a type that defines Debugf.
type DebugfReceiver interface {
	Debugf(string, ...interface{})
	DebugfContext(context.Context, string, ...interface{})
}

// OutputReceiver is an interface
type OutputReceiver interface {
	InfoReceiver
	InfofReceiver
	DebugReceiver
	DebugfReceiver
}

// WarningfReceiver is a type that defines Warningf.
type WarningfReceiver interface {
	Warningf(string, ...interface{})
	WarningfContext(context.Context, string, ...interface{})
}

// ErrorfReceiver is a type that defines Errorf.
type ErrorfReceiver interface {
	Errorf(string, ...interface{})
	ErrorfContext(context.Context, string, ...interface{})
}

// FatalfReceiver is a type that defines Fatalf.
type FatalfReceiver interface {
	Fatalf(string, ...interface{})
	FatalfContext(context.Context, string, ...interface{})
}

// ErrorOutputReceiver is an interface
type ErrorOutputReceiver interface {
	WarningfReceiver
	ErrorfReceiver
	FatalfReceiver
}

// WarningReceiver is a type that defines Warning.
type WarningReceiver interface {
	Warning(error, ...ErrorEventOption)
	WarningContext(context.Context, error, ...ErrorEventOption)
}

// ErrorReceiver is a type that defines Error.
type ErrorReceiver interface {
	Error(error, ...ErrorEventOption)
	ErrorContext(context.Context, error, ...ErrorEventOption)
}

// FatalReceiver is a type that defines Fatal.
type FatalReceiver interface {
	Fatal(error, ...ErrorEventOption)
	FatalContext(context.Context, error, ...ErrorEventOption)
}

// Errorable is an interface
type Errorable interface {
	WarningReceiver
	ErrorReceiver
	FatalReceiver
}

// Closer is a type that can close.
type Closer interface {
	Close()
}

// Drainer is a type that can be drained
type Drainer interface {
	Drain()
}

// FatalCloser is a type that defines Fatal and Close.
type FatalCloser interface {
	FatalReceiver
	Closer
}

// Flagged is a type that returns flags.
type Flagged interface {
	GetFlags() *Flags
	GetWritable() *Flags
}

// Log is a logger that implements the full suite of logging methods.
type Log interface {
	Scoper
	Triggerable
	OutputReceiver
	ErrorOutputReceiver
	Errorable
}

// FullLog is a logger that implements the full suite of logging methods.
type FullLog interface {
	Closer
	Drainer
	Flagged
	Listenable
	Filterable
	Log
}

package logger

import "context"

// Listenable is an interface loggers can ascribe to.
type Listenable interface {
	Listen(flag string, label string, listener Listener)
}

// Triggerable is type that can trigger events.
type Triggerable interface {
	Trigger(context.Context, Event)
}

// Scoper is a type that can return a scope.
type Scoper interface {
	// Apply augments a given context with fields from the Scope, including Labels, Annotations, and Path.
	Apply(context.Context) context.Context
	// WithContext sets the default context for the scope, which otherwise would be `context.Background()`
	WithContext(context.Context) Scope
	// WithPath returns a new scope with a given set of additional path segments.
	WithPath(...string) Scope
	// WithLabels returns a new scope with a given set of additional label values.
	WithLabels(Labels) Scope
	// WithAnnotations returns a new scope with a given set of additional annotation values.
	WithAnnotations(Annotations) Scope
}

// Writable is an type that can write events.
type Writable interface {
	Write(context.Context, Event)
}

// WriteTriggerable is a type that can both trigger and write events.
type WriteTriggerable interface {
	Triggerable
	Writable
}

// InfoReceiver is a type that defines Info.
type InfoReceiver interface {
	Info(...interface{})
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
}

// DebugReceiver is a type that defines Debug.
type DebugReceiver interface {
	Debug(...interface{})
}

// DebugfReceiver is a type that defines Debugf.
type DebugfReceiver interface {
	Debugf(string, ...interface{})
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
}

// ErrorfReceiver is a type that defines Errorf.
type ErrorfReceiver interface {
	Errorf(string, ...interface{})
}

// FatalfReceiver is a type that defines Fatalf.
type FatalfReceiver interface {
	Fatalf(string, ...interface{})
}

// ErrorOutputReceiver is an interface
type ErrorOutputReceiver interface {
	WarningfReceiver
	ErrorfReceiver
	FatalfReceiver
}

// WarningReceiver is a type that defines Warning.
type WarningReceiver interface {
	Warning(error, ...ErrorEventOption) error
}

// ErrorReceiver is a type that defines Error.
type ErrorReceiver interface {
	Error(error, ...ErrorEventOption) error
}

// FatalReceiver is a type that defines Fatal.
type FatalReceiver interface {
	Fatal(error, ...ErrorEventOption) error
}

// Errorable is an interface
type Errorable interface {
	WarningReceiver
	ErrorReceiver
	FatalReceiver
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
	Listenable
	Log
}

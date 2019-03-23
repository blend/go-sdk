package logger

import (
	"io"
	"time"
)

// Event is an interface representing methods necessary to trigger listeners.
type Event interface {
	Flag() string
	Timestamp() time.Time
}

// Listener is a function that can be triggered by events.
type Listener func(e Event)

// Listenable is an interface.
type Listenable interface {
	Listen(flag string, label string, listener Listener)
}

// Triggerable is an interface.
type Triggerable interface {
	Trigger(Event)
}

// SyncTriggerable is an interface.
type SyncTriggerable interface {
	SyncTrigger(Event)
}

// InfofReceiver is a type that defines Infof.
type InfofReceiver interface {
	Infof(string, ...interface{})
}

// SillyfReceiver is a type that defines Sillyf.
type SillyfReceiver interface {
	Sillyf(string, ...interface{})
}

// DebugfReceiver is a type that defines Debugf.
type DebugfReceiver interface {
	Debugf(string, ...interface{})
}

// OutputReceiver is an interface
type OutputReceiver interface {
	InfofReceiver
	SillyfReceiver
	DebugfReceiver
}

// SyncInfofReceiver is a type that defines SyncInfof.
type SyncInfofReceiver interface {
	SyncInfof(string, ...interface{})
}

// SyncSillyfReceiver is a type that defines SyncSillyf.
type SyncSillyfReceiver interface {
	SyncSillyf(string, ...interface{})
}

// SyncDebugfReceiver is a type that defines SyncDebugf.
type SyncDebugfReceiver interface {
	SyncDebugf(string, ...interface{})
}

// SyncOutputReceiver is an interface
type SyncOutputReceiver interface {
	SyncInfofReceiver
	SyncSillyfReceiver
	SyncDebugfReceiver
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

// SyncWarningfReceiver is a type that defines SyncWarningf.
type SyncWarningfReceiver interface {
	SyncWarningf(string, ...interface{})
}

// SyncErrorffReceiver is a type that defines SyncErrorf.
type SyncErrorffReceiver interface {
	SyncErrorf(string, ...interface{})
}

// SyncFatalfReceiver is a type that defines SyncFatalf.
type SyncFatalfReceiver interface {
	SyncFatalf(string, ...interface{})
}

// SyncErrorOutputReceiver is an interface.
type SyncErrorOutputReceiver interface {
	SyncWarningfReceiver
	SyncErrorffReceiver
	SyncFatalfReceiver
}

// WarningReceiver is a type that defines Warning.
type WarningReceiver interface {
	Warning(error) error
}

// ErrorReceiver is a type that defines Error.
type ErrorReceiver interface {
	Error(error) error
}

// FatalReceiver is a type that defines Fatal.
type FatalReceiver interface {
	Fatal(error) error
}

// Errorable is an interface
type Errorable interface {
	WarningReceiver
	ErrorReceiver
	FatalReceiver
}

// SyncWarningReceiver is a type that defines SyncWarning.
type SyncWarningReceiver interface {
	SyncWarning(error) error
}

// SyncErrorReceiver is a type that defines SyncError.
type SyncErrorReceiver interface {
	SyncError(error) error
}

// SyncFatalReceiver is a type that defines SyncFatal.
type SyncFatalReceiver interface {
	SyncFatal(error) error
}

// SyncErrorable is an interface
type SyncErrorable interface {
	SyncWarningReceiver
	SyncErrorReceiver
	SyncFatalReceiver
}

// SyncLogger is a logger that implements syncronous methods.
type SyncLogger interface {
	Listenable
	SyncTriggerable
	SyncOutputReceiver
	SyncErrorOutputReceiver
	SyncErrorable
}

// AsyncLogger is a logger that implements async methods.
type AsyncLogger interface {
	Listenable
	Triggerable
	OutputReceiver
	ErrorOutputReceiver
	Errorable
}

// FullReceiver is every possible receiving / output interface.
type FullReceiver interface {
	SyncTriggerable
	SyncOutputReceiver
	SyncErrorOutputReceiver
	SyncErrorable
	Triggerable
	OutputReceiver
	ErrorOutputReceiver
	Errorable
}

// FullLogger is every possible interface, including listenable.
type FullLogger interface {
	Listenable
	SyncTriggerable
	SyncOutputReceiver
	SyncErrorOutputReceiver
	SyncErrorable
	Triggerable
	OutputReceiver
	ErrorOutputReceiver
	Errorable
}

// Log is an alias to full logger.
// It is speculative as useful.
type Log = FullLogger

// Writer is a type that can consume events.
type Writer interface {
	Write(Event) error
	WriteError(Event) error
	Output() io.Writer
	ErrorOutput() io.Writer
}

// --------------------------------------------------------------------------------
// testing helpers
// --------------------------------------------------------------------------------

// MarshalEvent marshals an object as a logger event.
func MarshalEvent(obj interface{}) (Event, bool) {
	typed, isTyped := obj.(Event)
	return typed, isTyped
}

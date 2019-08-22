package logger

import "time"

// Flags
const (
	FlagAll  = "all"
	FlagNone = "none"
	Fatal    = "fatal"
	Error    = "error"
	Warning  = "warning"
	Debug    = "debug"
	Info     = "info"

	Audit = "audit"
)

// Output Formats
const (
	FormatJSON = "json"
	FormatText = "text"
)

// Default flags
var (
	DefaultFlags         = []string{Info, Error, Fatal}
	DefaultListenerName  = "default"
	DefaultRecoverPanics = true
)

// Environment Variable Names
const (
	EnvVarFlags      = "LOG_FLAGS"
	EnvVarFormat     = "LOG_FORMAT"
	EnvVarNoColor    = "NO_COLOR"
	EnvVarHideTime   = "LOG_HIDE_TIME"
	EnvVarTimeFormat = "LOG_TIME_FORMAT"
	EnvVarJSONPretty = "LOG_JSON_PRETTY"
)

const (
	// DefaultBufferPoolSize is the default buffer pool size.
	DefaultBufferPoolSize = 1 << 8 // 256
	// DefaultTextTimeFormat is the default time format.
	DefaultTextTimeFormat = time.RFC3339Nano
	// DefaultTextWriterUseColor is a default setting for writers.
	DefaultTextWriterUseColor = true
	// DefaultTextWriterShowHeadings is a default setting for writers.
	DefaultTextWriterShowHeadings = true
	// DefaultTextWriterShowTimestamp is a default setting for writers.
	DefaultTextWriterShowTimestamp = true
)

const (
	// DefaultWorkerQueueDepth is the default depth per listener to queue work.
	// It's currently set to 256k entries.
	DefaultWorkerQueueDepth = 1 << 10
)

// String constants
const (
	Space   = " "
	Newline = "\n"
)

// Common json fields
const (
	FieldFlag        = "flag"
	FieldTimestamp   = "_timestamp"
	FieldScopePath   = "scope_path"
	FieldText        = "text"
	FieldElapsed     = "elapsed"
	FieldLabels      = "labels"
	FieldAnnotations = "annotations"
)

// JSON Formatter defaults
const (
	DefaultJSONPretty = false
)

/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import "time"

const (
	// FlagAll enables all flags by default.
	FlagAll = "all"
	// FlagNone disables all flags by default.
	FlagNone = "none"
	// Fatal controls errors that should be considered process ending.
	Fatal = "fatal"
	// Error controls errors that should be logged by default and may affect user behavior.
	Error = "error"
	// Warning controls errors that should be skipped by default but may help debugging.
	Warning = "warning"
	// Debug controls output that is useful when diagnosing issues.
	Debug = "debug"
	// Info controls output that is useful for output by default.
	Info = "info"
	// Audit controls events that indiciate security related information.
	Audit = "audit"
)

const (
	// ScopeAll is a special scope that matches all scopes.
	ScopeAll = "*"
)

// Output Formats
const (
	FormatJSON = "json"
	FormatText = "text"
)

// Default flags
var (
	DefaultFlags          = []string{Info, Error, Fatal}
	DefaultFlagsWritable  = []string{FlagAll}
	DefaultScopes         = []string{ScopeAll}
	DefaultWritableScopes = []string{ScopeAll}
	DefaultListenerName   = "default"
	DefaultRecoverPanics  = true
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
	FieldRestricted  = "_restricted"
)

// JSON Formatter defaults
const (
	DefaultJSONPretty = false
)

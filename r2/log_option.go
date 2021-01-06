package r2

import "github.com/blend/go-sdk/sanitize"

// LogOptions are options that govern the logging of requests.
type LogOptions struct {
	RequestSanitizationDefaults []sanitize.RequestOption
}

// LogOption are mutators for log options.
type LogOption func(*LogOptions)

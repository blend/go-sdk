package datadog

const (
	// DefaultPort is the default port.
	DefaultPort = "8125"
	// DefaultTracePort is the default trace port.
	DefaultTracePort = "8126"
	// DefaultTracingEnabled is the default value for tracing enabled.
	DefaultTracingEnabled = true
	// DefaultAddress is the default address for datadog.
	DefaultAddress = "unix:///var/run/datadog/dsd.socket"
)

// Default Tags
const (
	TagService  = "service"
	TagEnv      = "env"
	TagHostname = "hostname"
)

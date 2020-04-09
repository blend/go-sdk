package datadog

import (
	"context"

	"github.com/blend/go-sdk/configutil"
)

const (
	// DefaultDatadogBufferDepth is the default number of statsd messages to buffer.
	DefaultDatadogBufferDepth = 128
)

// Config is the datadog config.
type Config struct {
	// Address is the address of the datadog collector in the form of "hostname:port" or "unix:///path/to/socket"
	Address string `json:"address,omitempty" yaml:"address,omitempty" env:"DATADOG_ADDRESS"`
	// Address is the address of the datadog collector in the form of "hostname:port" or "unix:///path/to/trace-socket"
	TraceAddress string `json:"traceAddress,omitempty" yaml:"traceAddress,omitempty" env:"DATADOG_TRACE_ADDRESS"`
	// TracingEnabled returns if we should use tracing or not.
	TracingEnabled *bool `json:"tracingEnabled" yaml:"tracingEnabled" env:"DATADOG_APM_ENABLED"`
	// Buffered indicates if we should buffer statsd messages or not.
	Buffered *bool `json:"buffered,omitempty" yaml:"buffered,omitempty" env:"DATADOG_BUFFERED"`
	// BufferDepth is the depth of the buffer for datadog events.
	// A zero value implies an unbuffered client.
	BufferDepth int `json:"bufferDepth,omitempty" yaml:"bufferDepth,omitempty" env:"DATADOG_BUFFER_DEPTH"`
	// Namespace is an optional namespace.
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty" env:"DATADOG_NAMESPACE"`
	// DefaultTags are the default tags associated with any stat metric.
	DefaultTags []string `json:"defaultTags,omitempty" yaml:"defaultTags,omitempty" env:"DATADOG_DEFAULT_TAGS,csv"`
}

// Resolve implements configutil.ConfigResolver.
func (c *Config) Resolve(ctx context.Context) error {
	return configutil.GetEnvVars(ctx).ReadInto(c)
}

// IsZero returns if the config is unset.
func (c Config) IsZero() bool {
	return c.Address == "" && c.TraceAddress == ""
}

// TracingEnabledOrDefault returns if tracing is enabled.
func (c Config) TracingEnabledOrDefault() bool {
	if c.TracingEnabled != nil {
		return *c.TracingEnabled
	}
	return DefaultTracingEnabled
}

// GetAddress returns the datadog collector address string.
func (c Config) GetAddress() string {
	if c.Address != "" {
		return c.Address
	}
	return DefaultAddress
}

// GetTraceAddress returns the datadog collector address string.
func (c Config) GetTraceAddress() string {
	if c.TraceAddress != "" {
		return c.TraceAddress
	}
	return DefaultTraceAddress
}

// BufferedOrDefault returns if the client should buffer messages or not.
func (c Config) BufferedOrDefault() bool {
	if c.Buffered != nil {
		return *c.Buffered
	}
	return false
}

// BufferDepthOrDefault returns the buffer depth.
func (c Config) BufferDepthOrDefault() int {
	if c.BufferDepth > 0 {
		return c.BufferDepth
	}
	return DefaultDatadogBufferDepth
}

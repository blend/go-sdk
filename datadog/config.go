/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package datadog

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/env"
)

const (
	// DefaultDatadogBufferDepth is the default number of statsd messages to buffer.
	DefaultDatadogBufferDepth = 128
)

// Config is the datadog config.
type Config struct {
	// Address is the address of the datadog collector in the form of "hostname:port" or "unix:///path/to/socket"
	// It will supercede `Hostname` and `Port`.
	Address	string	`json:"address,omitempty" yaml:"address,omitempty" env:"DATADOG_ADDRESS"`
	// TraceAddress is the address of the datadog collector in the form of "hostname:port" or "unix:///path/to/trace-socket"
	// It will supercede `TraceHostname` and `TracePort`
	TraceAddress	string	`json:"traceAddress,omitempty" yaml:"traceAddress,omitempty" env:"DATADOG_TRACE_ADDRESS"`
	// ProfilerAddress is the address of the datadog collector in the form of "hostname:port" or "unix:///path/to/profiler-socket"
	ProfilerAddress	string	`json:"profilerAddress,omitempty" yaml:"profilerAddress,omitempty" env:"DATADOG_PROFILER_ADDRESS"`

	// Hostname is the host portion of a <host>:<port> address. It will be used in conjunction with `Port`
	// to form the default `Address`.
	Hostname	string	`json:"hostname,omitempty" yaml:"hostname,omitempty" env:"DATADOG_HOSTNAME"`

	// CentralHostname is meant to replace Hostname when a central server is used for metrics collection.
	// CentralHostname is the host portion of <host>:<port> address. It will be used in conjunction with `Port`
	// to form the default `Address`.
	CentralHostname	string	`json:"centralHostname,omitempty" yaml:"centralHostname,omitempty" env:"DATADOG_CENTRAL_HOSTNAME"`
	// Port is the port portion of a <host>:<port> address. It will be used in conjunction with `Host`
	// to form the default `Address`.
	Port	string	`json:"port,omitempty" yaml:"port,omitempty" env:"DATADOG_PORT"`

	// TraceHostname is the host portion of a <host>:<port> address. It will be used in conjunction with `TracePort`
	// to form the default `TraceAddress`.
	TraceHostname	string	`json:"traceHostname,omitempty" yaml:"traceHostname,omitempty" env:"DATADOG_TRACE_HOSTNAME"`
	// TracePort is the port portion of a <host>:<port> address. It will be used in conjunction with `TraceHost`
	// to form the default `TraceAddress`.
	TracePort	string	`json:"tracePort,omitempty" yaml:"tracePort,omitempty" env:"DATADOG_TRACE_PORT"`

	// TracingEnabled returns if we should use tracing or not.
	TracingEnabled	*bool	`json:"tracingEnabled,omitempty" yaml:"tracingEnabled,omitempty" env:"DATADOG_APM_ENABLED"`
	// TracingSampleRate is the default tracing sample rate, on the interval [0-1]
	TraceSampleRate	*float64	`json:"traceSampleRate,omitempty" yaml:"traceSampleRate,omitempty" env:"DATADOG_APM_SAMPLE_RATE"`

	// ProfilingEnabled returns if we should use profiling or not.
	ProfilingEnabled	*bool	`json:"profilingEnabled,omitempty" yaml:"profilingEnabled,omitempty" env:"DATADOG_PROFILING_ENABLED"`

	// Buffered indicates if we should buffer statsd metrics.
	Buffered	*bool	`json:"buffered,omitempty" yaml:"buffered,omitempty" env:"DATADOG_BUFFERED"`
	// BufferDepth is the depth of the buffer for statsd metrics.
	BufferDepth	int	`json:"bufferDepth,omitempty" yaml:"bufferDepth,omitempty" env:"DATADOG_BUFFER_DEPTH"`

	// Namespace is an optional namespace.
	// The namespace is a prefix on all statsd metric names submitted to the collector.
	Namespace	string	`json:"namespace,omitempty" yaml:"namespace,omitempty" env:"DATADOG_NAMESPACE"`
	// DefaultTags are the default tags associated with any stat metric.
	DefaultTags	[]string	`json:"defaultTags,omitempty" yaml:"defaultTags,omitempty" env:"DATADOG_DEFAULT_TAGS,csv"`
}

// Resolve implements configutil.ConfigResolver.
func (c *Config) Resolve(ctx context.Context) error {
	return env.GetVars(ctx).ReadInto(c)
}

// IsZero returns if the config is unset.
func (c Config) IsZero() bool {
	return c.GetAddress() == ""
}

// GetAddress returns the datadog collector address string.
func (c Config) GetAddress() string {
	if c.CentralHostname != "" {
		return fmt.Sprintf("%s:%s", c.CentralHostname, c.PortOrDefault())
	}
	if c.Address != "" {
		return c.Address
	}
	if c.Hostname != "" {
		return fmt.Sprintf("%s:%s", c.Hostname, c.PortOrDefault())
	}
	return DefaultAddress
}

// GetTraceAddress returns the datadog collector address string.
func (c Config) GetTraceAddress() string {
	if c.TraceAddress != "" {
		return c.TraceAddress
	}
	if c.TraceHostname != "" {
		return fmt.Sprintf("%s:%s", c.TraceHostname, c.TracePortOrDefault())
	}
	if c.Hostname != "" {
		return fmt.Sprintf("%s:%s", c.Hostname, c.TracePortOrDefault())
	}
	return ""
}

// GetProfilerAddress gets the profiler address.
func (c Config) GetProfilerAddress() string {
	if c.ProfilerAddress != "" {
		return c.ProfilerAddress
	}
	return c.GetTraceAddress()
}

// PortOrDefault returns the port or a default.
func (c Config) PortOrDefault() string {
	if c.Port != "" {
		return c.Port
	}
	return DefaultPort
}

// TracePortOrDefault returns the trace port or a default.
func (c Config) TracePortOrDefault() string {
	if c.TracePort != "" {
		return c.TracePort
	}
	return DefaultTracePort
}

// TracingEnabledOrDefault returns if tracing is enabled.
func (c Config) TracingEnabledOrDefault() bool {
	if c.TracingEnabled != nil {
		return *c.TracingEnabled
	}
	return DefaultTracingEnabled
}

// TraceSampleRateOrDefault returns the trace sample rate or a default.
func (c Config) TraceSampleRateOrDefault() float64 {
	if c.TraceSampleRate != nil {
		return *c.TraceSampleRate
	}
	return DefaultTraceSampleRate
}

// ProfilingEnabledOrDefault returns if tracing is enabled.
func (c Config) ProfilingEnabledOrDefault() bool {
	if c.ProfilingEnabled != nil {
		return *c.ProfilingEnabled
	}
	return DefaultProfilingEnabled
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

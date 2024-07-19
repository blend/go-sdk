/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package datadog

const (
	// DefaultPort is the default port.
	DefaultPort = "8125"
	// DefaultTracePort is the default trace port.
	DefaultTracePort = "8126"
	// DefaultTracingEnabled is the default value for tracing enabled.
	DefaultTracingEnabled = true
	// DefaultProfilingEnabled is the default value for profiling enabled.
	DefaultProfilingEnabled = true
	// DefaultTraceSampleRate returns the default trace sample rate of 25%
	DefaultTraceSampleRate = 0.25
	// DefaultAddress is the default address for datadog.
	DefaultAddress = "unix:///var/run/datadog/dsd.socket"
)

// Default Tags
const (
	TagService  = "service"
	TagEnv      = "env"
	TagHostname = "hostname"
)

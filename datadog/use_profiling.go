/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package datadog

// UseProfiler returns if profiling is enabled and the profiler address is configured.
//
// It should be used to gate if you should create a profiler with `profiler.Start`.
func UseProfiler(cfg Config) bool {
	return cfg.ProfilingEnabledOrDefault() && cfg.GetProfilerAddress() != ""
}

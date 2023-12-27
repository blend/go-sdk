/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package datadog

// UseProfiler returns if profiling is enabled and the profiler address is configured.
//
// It should be used to gate if you should create a profiler with `profiler.Start`.
func UseProfiler(cfg Config) bool {
	return cfg.ProfilingEnabledOrDefault() && cfg.GetProfilerAddress() != ""
}

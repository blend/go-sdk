/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package datadog

// UseTracing returns if tracing is enabled and the trace address is configured.
//
// It should be used to gate if you should create tracers with `NewTracer`.
func UseTracing(cfg Config) bool {
	return cfg.TracingEnabledOrDefault() && cfg.GetTraceAddress() != ""
}

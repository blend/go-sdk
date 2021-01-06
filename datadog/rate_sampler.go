package datadog

import (
	"math/rand"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	_ ddtracer.RateSampler = (*RateSampler)(nil)
)

// RateSampler samples from a sample rate.
type RateSampler float64

// SetRate is a no-op
func (r RateSampler) SetRate(newRate float64) {}

// Rate returns the rate.
func (r RateSampler) Rate() float64 {
	return float64(r)
}

// Sample returns true if the given span should be sampled.
func (r RateSampler) Sample(spn ddtrace.Span) bool {
	if r < 1 {
		return rand.Float64() < float64(r)
	}
	return true
}

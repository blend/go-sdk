package datadog

import (
	"fmt"
	"math/rand"

	"github.com/opentracing/opentracing-go"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/tracing"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	ddopentracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentracer"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// TraceAppNameWeb returns an app name.
func TraceAppNameWeb(serviceName string) string {
	return fmt.Sprintf("%s-web", serviceName)
}

// TraceAppNameHTTPClient returns an app name.
func TraceAppNameHTTPClient(serviceName string) string {
	return fmt.Sprintf("%s-http-client", serviceName)
}

// TraceAppNamePostgres returns an app name.
func TraceAppNamePostgres(serviceName string) string {
	return fmt.Sprintf("%s-postgres", serviceName)
}

// TraceAppNameVaultClient returns an app name.
func TraceAppNameVaultClient(serviceName string) string {
	return fmt.Sprintf("%s-vault-client", serviceName)
}

// TraceAppNameGRPC returns an app name.
func TraceAppNameGRPC(serviceName string) string {
	return fmt.Sprintf("%s-grpc", serviceName)
}

// TraceAppNameGRPCClient returns an app name.
func TraceAppNameGRPCClient(serviceName string) string {
	return fmt.Sprintf("%s-grpc-client", serviceName)
}

// TraceAppNameCron returns an app name.
func TraceAppNameCron(serviceName string) string {
	return fmt.Sprintf("%s-cron", serviceName)
}

// OptTraceServiceEnv returns a dd tracer start option that sets the service environment.
func OptTraceServiceEnv(serviceEnv string) ddtracer.StartOption {
	return ddtracer.WithGlobalTag(tracing.TagKeyEnvironment, serviceEnv)
}

// OptTraceSampleRate returns a dd tracer start option that sets trace sampler with a given rate.
func OptTraceSampleRate(rate float64) ddtracer.StartOption {
	return ddtracer.WithSampler(RateSampler(rate))
}

// UseTracing returns if tracing is enabled and the trace address is configured.
//
// It should be used to gate if you should create tracers with `NewTracer`.
func UseTracing(cfg Config) bool {
	return cfg.TracingEnabledOrDefault() && cfg.GetTraceAddress() != ""
}

// NewTracer returns a new tracer.
//
// It defaults to the service environment as named by the `SERVICE_ENV` environment variable.
// It defaults the sample rate to the sample rate as returned by the configuration.
func NewTracer(appName string, cfg Config, opts ...ddtracer.StartOption) opentracing.Tracer {
	return ddopentracer.New(
		append([]ddtracer.StartOption{
			ddtracer.WithAgentAddr(cfg.GetTraceAddress()),
			ddtracer.WithServiceName(appName),
			OptTraceServiceEnv(env.Env().ServiceEnv()),
			OptTraceSampleRate(cfg.TraceSampleRateOrDefault()),
		}, opts...)...,
	)
}

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

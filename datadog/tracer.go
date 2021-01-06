package datadog

import (
	"context"
	"strings"

	"github.com/opentracing/opentracing-go"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/logger"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	ddopentracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentracer"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// NewTracer returns a new tracer.
//
// It defaults to the service environment as named by the `SERVICE_ENV` environment variable.
// It defaults the sample rate to the sample rate as returned by the configuration.
func NewTracer(opts ...TracerOption) opentracing.Tracer {
	options := TracerOptions{
		ServiceName: env.Env().ServiceName(),
		ServiceEnv:  env.Env().ServiceEnv(),
		Hostname:    env.Env().Hostname(),
		SampleRate:  DefaultTraceSampleRate,
	}
	for _, opt := range opts {
		opt(&options)
	}

	var startOptions []ddtracer.StartOption
	if options.Addr != "" {
		startOptions = append(startOptions, ddtracer.WithAgentAddr(options.Addr))
	}
	if options.StatsAddr != "" {
		startOptions = append(startOptions, ddtracer.WithAnalytics(true), ddtracer.WithAnalyticsRate(options.SampleRate))
		startOptions = append(startOptions, ddtracer.WithDogstatsdAddress(options.StatsAddr))
	}
	if options.ServiceEnv != "" {
		startOptions = append(startOptions, ddtracer.WithEnv(options.ServiceEnv))
	}
	if options.ServiceName != "" {
		startOptions = append(startOptions, ddtracer.WithService(options.ServiceName))
	}
	if options.Version != "" {
		startOptions = append(startOptions, ddtracer.WithServiceVersion(options.Version))
	}
	if options.Hostname != "" {
		startOptions = append(startOptions, ddtracer.WithGlobalTag(TagHostname, options.Hostname))
	}
	for key, value := range options.Tags {
		startOptions = append(startOptions, ddtracer.WithGlobalTag(key, value))
	}
	startOptions = append(startOptions, ddtracer.WithSampler(RateSampler(options.SampleRate)))
	if options.Log != nil {
		startOptions = append(startOptions, ddtracer.WithLogger(traceLogShim{options.Log}))
	} else {
		startOptions = append(startOptions, ddtracer.WithLogger(traceLogShim{}))
	}
	return ddopentracer.New(startOptions...)
}

// OptTraceAgentAddr returns a dd tracer start option that sets the agent addr.
func OptTraceAgentAddr(addr string) TracerOption {
	return func(to *TracerOptions) { to.Addr = addr }
}

// OptTraceServiceName returns a dd tracer start option that sets the service.
func OptTraceServiceName(serviceName string) TracerOption {
	return func(to *TracerOptions) { to.ServiceName = serviceName }
}

// OptTraceServiceEnv returns a dd tracer start option that sets the service environment.
func OptTraceServiceEnv(serviceEnv string) TracerOption {
	return func(to *TracerOptions) { to.ServiceEnv = serviceEnv }
}

// OptTraceVersion returns a dd tracer start option that sets the service version.
func OptTraceVersion(version string) TracerOption {
	return func(to *TracerOptions) { to.Version = version }
}

// OptTraceHostname returns a dd tracer start option that sets the service hostname.
func OptTraceHostname(hostname string) TracerOption {
	return func(to *TracerOptions) { to.Hostname = hostname }
}

// OptTraceSampleRate returns a dd tracer start option that sets trace sampler with a given rate.
func OptTraceSampleRate(rate float64) TracerOption {
	return func(to *TracerOptions) { to.SampleRate = rate }
}

// OptTraceLog returns an option that sets the log output.
func OptTraceLog(log logger.Triggerable) TracerOption {
	return func(to *TracerOptions) { to.Log = log }
}

// OptTraceConfig sets relevant fields from the datadog config.
func OptTraceConfig(cfg Config) TracerOption {
	return func(to *TracerOptions) {
		to.Addr = cfg.GetTraceAddress()
		to.StatsAddr = cfg.GetAddress()
		to.SampleRate = cfg.TraceSampleRateOrDefault()
	}
}

// TracerOption mutates tracer options
type TracerOption func(*TracerOptions)

// TracerOptions are all the options we can set when creating a tracer.
type TracerOptions struct {
	Addr        string
	StatsAddr   string
	ServiceName string
	ServiceEnv  string
	Version     string
	Hostname    string
	SampleRate  float64
	Tags        map[string]string
	Log         logger.Triggerable
}

var (
	_ ddtrace.Logger = (*traceLogShim)(nil)
)

// traceLogShim is a shim between the sdk logger and the datadog logger.
type traceLogShim struct {
	Logger logger.Triggerable
}

// Log implements ddtrace.Logger
func (tls traceLogShim) Log(msg string) {
	if tls.Logger != nil {
		tls.Logger.TriggerContext(context.Background(), logger.NewMessageEvent("datadog-tracer", strings.TrimSpace(msg)))
	}
}

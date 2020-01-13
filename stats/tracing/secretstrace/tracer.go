package secretstrace

import (
	"context"
	"strconv"
	"time"
	"github.com/blend/go-sdk/secrets"
	"github.com/blend/go-sdk/stats/tracing"
	"github.com/opentracing/opentracing-go"
)

// Tracer returns a request tracer that also injects span context into outgoing headers.
func Tracer(tracer opentracing.Tracer) secrets.Tracer {
	return &secretsTracer{tracer: tracer}
}

type secretsTracer struct {
	tracer opentracing.Tracer
}

func (st secretsTracer) Start(ctx context.Context, options ...secrets.TraceOption) (secrets.TraceFinisher, error) {
	var config secrets.SecretTraceConfig
	for _, opt := range options {
		err := opt(&config)
		if err != nil {
			return secretsTraceFinisher{}, nil
		}
	}
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeVault},
		opentracing.StartTime(time.Now().UTC()),
	}
	if config.VaultOperation != "" {
		startOptions = append(startOptions, opentracing.Tag{Key: tracing.TagSecretsOperation, Value: config.VaultOperation})
	}
	if config.KeyName != "" {
		startOptions = append(startOptions, opentracing.Tag{Key: tracing.TagSecretKey, Value: config.KeyName})
	}
	span, _ := tracing.StartSpanFromContext(ctx, st.tracer, tracing.OperationVaultAPI, startOptions...)
	return secretsTraceFinisher{span: span}, nil
}

type secretsTraceFinisher struct {
	span opentracing.Span
}

func (stf secretsTraceFinisher) Finish(_ context.Context, vaultStatusCode int, vaultError error) {
	if stf.span == nil {
		return
	}
	if vaultError != nil {
		tracing.SpanError(stf.span, vaultError)
	}
	stf.span.SetTag(tracing.TagKeyHTTPCode, strconv.Itoa(vaultStatusCode))
	stf.span.Finish()
}

/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package vaulttrace

import (
	"context"
	"strconv"
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/blend/go-sdk/tracing"
	"github.com/blend/go-sdk/vault"
)

// Tracer returns a request tracer that also injects span context into outgoing headers.
func Tracer(tracer opentracing.Tracer) vault.Tracer {
	return &vaultTracer{tracer: tracer}
}

type vaultTracer struct {
	tracer opentracing.Tracer
}

func (vt vaultTracer) Start(ctx context.Context, options ...vault.TraceOption) (vault.TraceFinisher, error) {
	var config vault.SecretTraceConfig
	for _, opt := range options {
		err := opt(&config)
		if err != nil {
			return vaultTraceFinisher{}, nil
		}
	}
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeVault},
		tracing.TagMeasured(),
		opentracing.StartTime(time.Now().UTC()),
	}
	if config.VaultOperation != "" {
		startOptions = append(startOptions, opentracing.Tag{Key: tracing.TagSecretsOperation, Value: config.VaultOperation})
	}
	if config.KeyName != "" {
		startOptions = append(startOptions, opentracing.Tag{Key: tracing.TagSecretKey, Value: config.KeyName})
	}
	span, _ := tracing.StartSpanFromContext(ctx, vt.tracer, tracing.OperationVaultAPI, startOptions...)
	return vaultTraceFinisher{span: span}, nil
}

type vaultTraceFinisher struct {
	span opentracing.Span
}

func (vtf vaultTraceFinisher) Finish(_ context.Context, vaultStatusCode int, vaultError error) {
	if vtf.span == nil {
		return
	}
	if vaultError != nil {
		tracing.SpanError(vtf.span, vaultError)
	}
	vtf.span.SetTag(tracing.TagKeyHTTPCode, strconv.Itoa(vaultStatusCode))
	vtf.span.Finish()
}

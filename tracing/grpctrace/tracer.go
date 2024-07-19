/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package grpctrace

import (
	"context"
	"time"

	opentracing "github.com/opentracing/opentracing-go"

	"google.golang.org/grpc/metadata"

	"github.com/blend/go-sdk/grpcutil"
	"github.com/blend/go-sdk/tracing"
)

var (
	_ grpcutil.ClientTracer = (*tracer)(nil)
	_ grpcutil.ServerTracer = (*tracer)(nil)
)

// Tracer returns a tracer.
func Tracer(t opentracing.Tracer) grpcutil.Tracer {
	return tracer{
		Tracer: t,
	}
}

// Tracer implements grpcutil.ClientTracer and grpcutil.ServerTracer.
type tracer struct {
	opentracing.Tracer
}

// StartClientUnary starts a unary client trace.
func (t tracer) StartClientUnary(ctx context.Context, remoteAddr, method string) (context.Context, grpcutil.TraceFinisher, error) {
	finisher := TraceFinisher{
		startTime: time.Now().UTC(),
	}

	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeRPC},
		opentracing.Tag{Key: tracing.TagKeyResourceName, Value: method},
		opentracing.Tag{Key: tracing.TagKeyGRPCRemoteAddr, Value: remoteAddr},
		opentracing.Tag{Key: tracing.TagKeyGRPCMethod, Value: method},
		opentracing.Tag{Key: tracing.TagKeyGRPCRole, Value: "client"},
		opentracing.Tag{Key: tracing.TagKeyGRPCCallingConvention, Value: "unary"},
		tracing.TagMeasured(),
		opentracing.StartTime(finisher.startTime),
	}

	finisher.span, finisher.ctx = tracing.StartSpanFromContext(ctx, t.Tracer, tracing.OperationGRPCClientUnary, startOptions...)
	md := make(metadata.MD)
	if err := t.Tracer.Inject(finisher.span.Context(), opentracing.TextMap, MetadataReaderWriter{md}); err != nil {
		return nil, nil, err
	}
	finisher.ctx = metadata.NewOutgoingContext(finisher.ctx, md)
	return finisher.ctx, finisher, nil
}

// StartClientStream starts a stream client trace.
func (t tracer) StartClientStream(ctx context.Context, remoteAddr, method string) (context.Context, grpcutil.TraceFinisher, error) {
	var finisher TraceFinisher
	finisher.ctx = ctx
	finisher.startTime = time.Now().UTC()

	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeRPC},
		opentracing.Tag{Key: tracing.TagKeyGRPCRemoteAddr, Value: remoteAddr},
		opentracing.Tag{Key: tracing.TagKeyResourceName, Value: method},
		opentracing.Tag{Key: tracing.TagKeyGRPCMethod, Value: method},
		opentracing.Tag{Key: tracing.TagKeyGRPCRole, Value: "client"},
		opentracing.Tag{Key: tracing.TagKeyGRPCCallingConvention, Value: "stream"},
		opentracing.StartTime(finisher.startTime),
	}

	finisher.span, finisher.ctx = tracing.StartSpanFromContext(finisher.ctx, t.Tracer, tracing.OperationGRPCClientStream, startOptions...)
	md := make(metadata.MD)
	if err := t.Tracer.Inject(finisher.span.Context(), opentracing.TextMap, MetadataReaderWriter{md}); err != nil {
		return nil, nil, err
	}
	finisher.ctx = metadata.NewOutgoingContext(finisher.ctx, md)
	return finisher.ctx, finisher, nil
}

// StartServerUnary starts a unary server trace.
func (t tracer) StartServerUnary(ctx context.Context, method string) (context.Context, grpcutil.TraceFinisher, error) {
	var finisher TraceFinisher
	finisher.startTime = time.Now().UTC()
	finisher.ctx = ctx

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	authority := grpcutil.MetaValue(md, grpcutil.MetaTagAuthority)
	contentType := grpcutil.MetaValue(md, grpcutil.MetaTagContentType)
	userAgent := grpcutil.MetaValue(md, grpcutil.MetaTagUserAgent)

	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeRPC},
		opentracing.Tag{Key: tracing.TagKeyResourceName, Value: method},
		opentracing.Tag{Key: tracing.TagKeyGRPCMethod, Value: method},
		opentracing.Tag{Key: tracing.TagKeyGRPCAuthority, Value: authority},
		opentracing.Tag{Key: tracing.TagKeyGRPCUserAgent, Value: userAgent},
		opentracing.Tag{Key: tracing.TagKeyGRPCContentType, Value: contentType},
		opentracing.Tag{Key: tracing.TagKeyGRPCRole, Value: "server"},
		opentracing.Tag{Key: tracing.TagKeyGRPCCallingConvention, Value: "unary"},
		opentracing.StartTime(finisher.startTime),
	}

	// try to extract an incoming span context
	// this is typically done if we're a service being called in a chain from another (more ancestral)
	// span context.
	spanContext, _ := t.Tracer.Extract(opentracing.HTTPHeaders, MetadataReaderWriter{md})
	if spanContext != nil {
		startOptions = append(startOptions, opentracing.ChildOf(spanContext))
	}

	finisher.span = t.Tracer.StartSpan(tracing.OperationGRPCServerUnary, startOptions...)
	finisher.ctx = opentracing.ContextWithSpan(finisher.ctx, finisher.span)
	return finisher.ctx, finisher, nil
}

// StartServerStream starts a stream server trace.
func (t tracer) StartServerStream(ctx context.Context, method string) (context.Context, grpcutil.TraceFinisher, error) {
	var finisher TraceFinisher
	finisher.startTime = time.Now().UTC()
	finisher.ctx = ctx

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	authority := grpcutil.MetaValue(md, grpcutil.MetaTagAuthority)
	contentType := grpcutil.MetaValue(md, grpcutil.MetaTagContentType)
	userAgent := grpcutil.MetaValue(md, grpcutil.MetaTagUserAgent)
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeRPC},
		opentracing.Tag{Key: tracing.TagKeyResourceName, Value: method},
		opentracing.Tag{Key: tracing.TagKeyGRPCMethod, Value: method},
		opentracing.Tag{Key: tracing.TagKeyGRPCAuthority, Value: authority},
		opentracing.Tag{Key: tracing.TagKeyGRPCUserAgent, Value: userAgent},
		opentracing.Tag{Key: tracing.TagKeyGRPCContentType, Value: contentType},
		opentracing.Tag{Key: tracing.TagKeyGRPCRole, Value: "server"},
		opentracing.Tag{Key: tracing.TagKeyGRPCCallingConvention, Value: "stream"},
		opentracing.StartTime(finisher.startTime),
	}

	// try to extract an incoming span context
	// this is typically done if we're a service being called in a chain from another (more ancestral)
	// span context.
	spanContext, _ := t.Tracer.Extract(opentracing.HTTPHeaders, MetadataReaderWriter{md})
	if spanContext != nil {
		startOptions = append(startOptions, opentracing.ChildOf(spanContext))
	}
	finisher.span = t.Tracer.StartSpan(tracing.OperationGRPCServerStream, startOptions...)
	finisher.ctx = opentracing.ContextWithSpan(ctx, finisher.span)
	return ctx, finisher, nil
}

// TraceFinisher finishes traces.
type TraceFinisher struct {
	startTime time.Time
	ctx       context.Context
	span      opentracing.Span
}

// Finish implements TraceFinisher.
func (tf TraceFinisher) Finish(err error) {
	if err != nil {
		tracing.SpanError(tf.span, err)
	}
	tf.span.Finish()
}

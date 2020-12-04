package tracing

// SpanIDProvider is a tracing span context that has a SpanID getter
type SpanIDProvider interface {
	SpanID() uint64
}

// TraceIDProvider is a tracing span context that has a TraceID getter
type TraceIDProvider interface {
	TraceID() uint64
}

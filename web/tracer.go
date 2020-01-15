package web

// Tracer is a type that traces complete requests.
type Tracer interface {
	Start(*Ctx) TraceFinisher
}

// TraceFinisher is a finisher for a trace.
type TraceFinisher interface {
	Finish(*Ctx, error)
}

// ViewTracer is a type that can listen for view rendering traces.
type ViewTracer interface {
	StartView(*Ctx, *ViewResult) ViewTraceFinisher
}

// ViewTraceFinisher is a finisher for view traces.
type ViewTraceFinisher interface {
	FinishView(*Ctx, *ViewResult, error)
}

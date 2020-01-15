package web

import (
	"net/http"
)

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

// HTTPTracer is a simplified version of `Tracer` intended for a raw
// `(net/http).Request`. It returns a "new" request the request context may
// be modified after opening a span.
type HTTPTracer interface {
	Start(*http.Request) (HTTPTraceFinisher, *http.Request)
}

// HTTPTraceFinisher is a simplified version of `TraceFinisher` intended for a
// raw `(net/http).Request`.
type HTTPTraceFinisher interface {
	Finish(*http.Request, error)
}

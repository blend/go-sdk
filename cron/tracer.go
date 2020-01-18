package cron

import "context"

// Tracer is a trace handler.
type Tracer interface {
	Start(context.Context) (context.Context, TraceFinisher)
}

// TraceFinisher is a finisher for traces.
type TraceFinisher interface {
	Finish(context.Context, error)
}

package oauth

import "golang.org/x/oauth2"

// Tracer is a trace shim.
type Tracer interface {
	Start(r *oauth2.Config) TraceFinisher
}

// TraceFinisher is a finisher for a trace.
type TraceFinisher interface {
	Finish(*oauth2.Config, *Result, error)
}

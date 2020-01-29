package webutil

import (
	"net/http"
)

// TagKV is a key value pair for a span tag.
type TagKV struct {
	Key   string
	Value interface{}
}

// HTTPTracer is a simplified version of `web.Tracer` intended for a raw
// `(net/http).Request`. It returns a "new" request because the request context
// may be modified after opening a span.
type HTTPTracer interface {
	Start(*http.Request, ...TagKV) (HTTPTraceFinisher, *http.Request)
}

// HTTPTraceFinisher is a finisher for `HTTPTracer`.
type HTTPTraceFinisher interface {
	Finish(error)
}

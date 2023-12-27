/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package traceserver

type (
	// SpanList implements msgp.Encodable on top of a slice of spans.
	SpanList	[]*Span

	// SpanLists implements msgp.Decodable on top of a slice of spanList.
	// This type is only used in tests.
	SpanLists	[]SpanList
)

// Span represents a computation. Callers must call Finish when a span is
// complete to ensure it's submitted.
type Span struct {
	Name		string			`json:"name" msg:"name"`				// operation name
	Service		string			`json:"service" msg:"service"`				// service name (i.e. "grpc.server", "http.request")
	Resource	string			`json:"resource" msg:"resource"`			// resource name (i.e. "/user?id=123", "SELECT * FROM users")
	Type		string			`json:"type" msg:"type"`				// protocol associated with the span (i.e. "web", "db", "cache")
	Start		int64			`json:"start" msg:"start"`				// span start time expressed in nanoseconds since epoch
	Duration	int64			`json:"duration" msg:"duration"`			// duration of the span expressed in nanoseconds
	Meta		map[string]string	`json:"meta,omitempty" msg:"meta,omitempty"`		// arbitrary map of metadata
	Metrics		map[string]float64	`json:"metrics,omitempty" msg:"metrics,omitempty"`	// arbitrary map of numeric metrics
	SpanID		uint64			`json:"span_id" msg:"span_id"`				// identifier of this span
	TraceID		uint64			`json:"trace_id" msg:"trace_id"`			// identifier of the root span
	ParentID	uint64			`json:"parent_id" msg:"parent_id"`			// identifier of the span's direct parent
	Error		int32			`json:"error" msg:"error"`				// error status of the span; 0 means no errors
}

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package autoflush

import "context"

// Tracer is a type that can trace actions in the Buffer.
type Tracer interface {
	StartAdd(context.Context) TraceFinisher
	StartAddMany(context.Context) TraceFinisher
	StartQueueFlush(context.Context) TraceFinisher
	StartFlush(context.Context) (context.Context, TraceFinisher)
}

// TraceFinisher finishes traces.
type TraceFinisher interface {
	Finish(error)
}

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package cron

import "context"

// Tracer is a trace handler.
type Tracer interface {
	Start(context.Context, string) (context.Context, TraceFinisher)
}

// TraceFinisher is a finisher for traces.
type TraceFinisher interface {
	Finish(context.Context, error)
}

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package redis

import "context"

// Tracer is a type that can implement traces.
type Tracer interface {
	Do(context.Context, Config, string, []string) TraceFinisher
}

// TraceFinisher is a type that can finish traces.
type TraceFinisher interface {
	Finish(context.Context, error)
}

/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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

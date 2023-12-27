/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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

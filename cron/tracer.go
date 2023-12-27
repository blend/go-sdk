/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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

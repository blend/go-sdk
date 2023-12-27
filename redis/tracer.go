/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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

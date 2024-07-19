/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package db

import (
	"context"
	"database/sql"
)

// Tracer is a type that can implement traces.
// If any of the methods return a nil finisher, they will be skipped.
type Tracer interface {
	Prepare(context.Context, Config, string) TraceFinisher
	Query(context.Context, Config, string, string) TraceFinisher
}

// TraceFinisher is a type that can finish traces.
type TraceFinisher interface {
	FinishPrepare(context.Context, error)
	FinishQuery(context.Context, sql.Result, error)
}

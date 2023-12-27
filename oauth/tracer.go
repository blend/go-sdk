/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package oauth

import (
	"context"

	"golang.org/x/oauth2"
)

// Tracer is a trace shim.
type Tracer interface {
	Start(context.Context, *oauth2.Config) TraceFinisher
}

// TraceFinisher is a finisher for a trace.
type TraceFinisher interface {
	Finish(context.Context, *oauth2.Config, *Result, error)
}

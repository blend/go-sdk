/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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

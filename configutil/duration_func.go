/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

import (
	"context"
	"time"
)

var (
	_ DurationSource = (*DurationFunc)(nil)
)

// DurationFunc is a value source from a function.
type DurationFunc func(context.Context) (*time.Duration, error)

// Duration returns an invocation of the function.
func (vf DurationFunc) Duration(ctx context.Context) (*time.Duration, error) {
	return vf(ctx)
}

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

import "context"

var (
	_ Float64Source = (*Float64Func)(nil)
)

// Float64Func is a float value source from a commandline flag.
type Float64Func func(context.Context) (*float64, error)

// Float64 returns an invocation of the function.
func (vf Float64Func) Float64(ctx context.Context) (*float64, error) {
	return vf(ctx)
}

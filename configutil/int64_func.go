/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import "context"

var (
	_ Int64Source = (*Int64Func)(nil)
)

// Int64Func is an int64 value source from a commandline flag.
type Int64Func func(context.Context) (*int64, error)

// Int64 returns an invocation of the function.
func (vf Int64Func) Int64(ctx context.Context) (*int64, error) {
	return vf(ctx)
}

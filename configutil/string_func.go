/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import "context"

var (
	_ StringSource = (*StringFunc)(nil)
)

// StringFunc is a value source from a function.
type StringFunc func(context.Context) (*string, error)

// String returns an invocation of the function.
func (svf StringFunc) String(ctx context.Context) (*string, error) {
	return svf(ctx)
}

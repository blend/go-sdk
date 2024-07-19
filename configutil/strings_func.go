/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import "context"

var (
	_ StringsSource = (*StringsFunc)(nil)
)

// StringsFunc is a value source from a function.
type StringsFunc func(context.Context) ([]string, error)

// Strings returns an invocation of the function.
func (svf StringsFunc) Strings(ctx context.Context) ([]string, error) {
	return svf(ctx)
}

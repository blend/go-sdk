/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import "context"

// IntPtr returns an IntSource for a given int pointer.
func IntPtr(value *int) IntSource {
	return IntPtrSource{Value: value}
}

var (
	_ IntSource = (*IntPtrSource)(nil)
)

// IntPtrSource is a IntSource that wraps an int pointer.
type IntPtrSource struct {
	Value *int
}

// Int implements IntSource.
func (ips IntPtrSource) Int(_ context.Context) (*int, error) {
	return ips.Value, nil
}

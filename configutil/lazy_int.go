/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import "context"

// LazyInt returns an IntSource for a given int pointer.
//
// LazyInt differs from IntPtr in that it treats 0 values as unset.
// If 0 is a valid value, use a IntPtr.
func LazyInt(value *int) LazyIntSource {
	return LazyIntSource{Value: value}
}

var (
	_ IntSource = (*LazyIntSource)(nil)
)

// LazyIntSource implements value provider.
//
// Note: LazyInt treats 0 as unset, if 0 is a valid value you must use configutil.IntPtr.
type LazyIntSource struct {
	Value *int
}

// Int returns the value for a constant.
func (i LazyIntSource) Int(_ context.Context) (*int, error) {
	if i.Value != nil && *i.Value > 0 {
		return i.Value, nil
	}
	return nil, nil
}

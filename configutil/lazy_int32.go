/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

import "context"

// LazyInt32 returns an Int32Source for a given int32 pointer.
//
// LazyInt32 differs from Int32Ptr in that it treats 0 values as unset.
// If 0 is a valid value, use a Int32Ptr.
func LazyInt32(value *int32) LazyInt32Source {
	return LazyInt32Source{Value: value}
}

var (
	_ Int32Source = (*LazyInt32Source)(nil)
)

// LazyInt32Source implements value provider.
//
// Note: LazyInt32Source treats 0 as unset, if 0 is a valid value you must use configutil.Int32Ptr.
type LazyInt32Source struct {
	Value *int32
}

// Int32 returns the value for a constant.
func (i LazyInt32Source) Int32(_ context.Context) (*int32, error) {
	if i.Value != nil && *i.Value > 0 {
		return i.Value, nil
	}
	return nil, nil
}

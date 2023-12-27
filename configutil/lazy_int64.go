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

// LazyInt64 returns an Int64Source for a given int64 pointer.
//
// LazyInt64 differs from Int64Ptr in that it treats 0 values as unset.
// If 0 is a valid value, use a Int64Ptr.
func LazyInt64(value *int64) LazyInt64Source {
	return LazyInt64Source{Value: value}
}

var (
	_ Int64Source = (*LazyInt64Source)(nil)
)

// LazyInt64Source implements value provider.
//
// Note: LazyInt64Source treats 0 as unset, if 0 is a valid value you must use configutil.Int64Ptr.
type LazyInt64Source struct {
	Value *int64
}

// Int64 returns the value for a constant.
func (i LazyInt64Source) Int64(_ context.Context) (*int64, error) {
	if i.Value != nil && *i.Value > 0 {
		return i.Value, nil
	}
	return nil, nil
}

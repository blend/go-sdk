/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

import "context"

// LazyFloat64 returns an Float64Source for a given float64 pointer.
//
// LazyFloat64 differs from Float64Ptr in that it treats 0 values as unset.
// If 0 is a valid value, use a Float64Ptr.
func LazyFloat64(value *float64) LazyFloat64Source {
	return LazyFloat64Source{Value: value}
}

var (
	_ Float64Source = (*LazyFloat64Source)(nil)
)

// LazyFloat64Source implements value provider.
//
// Note: LazyFloat64Source treats 0 as unset, if 0 is a valid value you must use configutil.Float64Ptr.
type LazyFloat64Source struct {
	Value *float64
}

// Float64 returns the value for a constant.
func (i LazyFloat64Source) Float64(_ context.Context) (*float64, error) {
	if i.Value != nil && *i.Value > 0 {
		return i.Value, nil
	}
	return nil, nil
}

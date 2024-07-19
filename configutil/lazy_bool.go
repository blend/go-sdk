/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import (
	"context"
)

// LazyBool returns an BoolSource for a given bool pointer.
//
// LazyBool differs from Bool in that it treats false values as unset.
// If false is a valid value, use a Bool.
func LazyBool(value *bool) LazyBoolSource {
	return LazyBoolSource{Value: value}
}

var (
	_ BoolSource = (*LazyBoolSource)(nil)
)

// LazyBoolSource implements value provider.
//
// Note: LazyDuration treats 0 as unset, if 0 is a valid value you must use configutil.DurationPtr.
type LazyBoolSource struct {
	Value *bool
}

// Bool returns the value for a constant.
func (i LazyBoolSource) Bool(_ context.Context) (*bool, error) {
	if i.Value != nil && *i.Value {
		return i.Value, nil
	}
	return nil, nil
}

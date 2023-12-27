/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

import "context"

var (
	_ IntSource = (*Int)(nil)
)

// Int implements value provider.
//
// Note: Int treats 0 as unset, if 0 is a valid value you must use configutil.IntPtr.
type Int int

// Int returns the value for a constant.
func (i Int) Int(_ context.Context) (*int, error) {
	if i > 0 {
		value := int(i)
		return &value, nil
	}
	return nil, nil
}

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

import "context"

var (
	_ Int32Source = (*Int32)(nil)
)

// Int32 implements value provider.
//
// Note: Int32 treats 0 as unset, if 0 is a valid value you must use configutil.Int32Ptr.
type Int32 int32

// Int32 returns the value for a constant.
func (i Int32) Int32(_ context.Context) (*int32, error) {
	if i > 0 {
		value := int32(i)
		return &value, nil
	}
	return nil, nil
}

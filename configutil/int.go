/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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

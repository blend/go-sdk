/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

import "context"

var (
	_ Float64Source = (*Float64)(nil)
)

// Float64 implements value provider.
//
// Note: this will treat 0 as unset, if 0 is a valid value you must use configutil.FloatPtr.
type Float64 float64

// Float64 returns the value for a constant.
func (f Float64) Float64(_ context.Context) (*float64, error) {
	if f > 0 {
		value := float64(f)
		return &value, nil
	}
	return nil, nil
}

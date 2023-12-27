/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

import (
	"context"
	"time"
)

var (
	_ DurationSource = (*Duration)(nil)
)

// Duration implements value provider.
//
// If the value is zero, a nil is returned by the implementation indicating
// the value was not present.
//
// If you want 0 to be a valid value, you must use DurationPtr.
type Duration time.Duration

// Duration returns the value for a constant.
func (dc Duration) Duration(_ context.Context) (*time.Duration, error) {
	if dc > 0 {
		value := time.Duration(dc)
		return &value, nil
	}
	return nil, nil
}

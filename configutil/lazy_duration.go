/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

import (
	"context"
	"time"
)

// LazyDuration returns an DurationSource for a given duration pointer.
//
// LazyDuration differs from DurationPtr in that it treats 0 values as unset.
// If 0 is a valid value, use a DurationPtr.
func LazyDuration(value *time.Duration) LazyDurationSource {
	return LazyDurationSource{Value: value}
}

var (
	_ DurationSource = (*LazyDurationSource)(nil)
)

// LazyDurationSource implements value provider.
//
// Note: LazyDuration treats 0 as unset, if 0 is a valid value you must use configutil.DurationPtr.
type LazyDurationSource struct {
	Value *time.Duration
}

// Duration returns the value for a constant.
func (i LazyDurationSource) Duration(_ context.Context) (*time.Duration, error) {
	if i.Value != nil && *i.Value > 0 {
		return i.Value, nil
	}
	return nil, nil
}

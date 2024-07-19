/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import (
	"context"
	"time"
)

// DurationPtr returns a DurationSource for a given duration pointer.
func DurationPtr(value *time.Duration) DurationSource {
	return DurationPtrSource{Value: value}
}

var (
	_ DurationSource = (*DurationPtrSource)(nil)
)

// DurationPtrSource is a DurationSource that wraps a duration pointer.
type DurationPtrSource struct {
	Value *time.Duration
}

// Duration implements DurationSource.
func (dps DurationPtrSource) Duration(_ context.Context) (*time.Duration, error) {
	return dps.Value, nil
}

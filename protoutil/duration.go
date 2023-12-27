/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package protoutil

import (
	"time"

	"github.com/golang/protobuf/ptypes/duration"
)

// Duration returns a proto duration.
func Duration(d time.Duration) *duration.Duration {
	if d == 0 {
		return nil
	}
	nanos := d.Nanoseconds()
	secs := nanos / 1e9
	nanos -= secs * 1e9

	return &duration.Duration{Seconds: int64(secs), Nanos: int32(nanos)}
}

// FromDuration returns a time.Duration.
func FromDuration(dur *duration.Duration) time.Duration {
	if dur == nil {
		return 0
	}
	d := time.Duration(dur.Seconds) * time.Second
	if dur.Nanos > 0 {
		d += time.Duration(dur.Nanos) * time.Nanosecond
	}
	return d
}

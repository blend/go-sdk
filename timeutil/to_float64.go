/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package timeutil

import "time"

// ToFloat64 returns a float64 representation of a time.
func ToFloat64(t time.Time) float64 {
	return float64(t.UnixNano())
}

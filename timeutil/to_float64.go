/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package timeutil

import "time"

// ToFloat64 returns a float64 representation of a time.
func ToFloat64(t time.Time) float64 {
	return float64(t.UnixNano())
}

/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cache

import "time"

// Stats represents cached statistics.
type Stats struct {
	Count		int
	SizeBytes	int
	MaxAge		time.Duration
}

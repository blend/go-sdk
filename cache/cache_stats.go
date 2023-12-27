/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package cache

import "time"

// Stats represents cached statistics.
type Stats struct {
	Count		int
	SizeBytes	int
	MaxAge		time.Duration
}

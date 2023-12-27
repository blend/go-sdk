/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package retry

import "time"

// Defaults
const (
	DefaultMaxAttempts	= 5
	DefaultRetryDelay	= time.Second
)

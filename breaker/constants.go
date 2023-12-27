/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package breaker

import "time"

// Constants
const (
	DefaultClosedExpiryInterval		= 5 * time.Second
	DefaultOpenExpiryInterval		= 60 * time.Second
	DefaultHalfOpenMaxActions	int64	= 1
	DefaultOpenFailureThreshold	int64	= 5
)

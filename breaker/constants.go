/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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

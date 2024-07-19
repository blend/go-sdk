/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package status

import (
	"time"

	"github.com/blend/go-sdk/ex"
)

// Errors
const (
	ErrServiceCheckNotDefined ex.Class = "service check is not defined for service"
)

const (
	// DefaultFreeformTimeout is a timeout.
	DefaultFreeformTimeout = 10 * time.Second
	// DefaultTrackedActionExpiration is the default tracker expiration.
	DefaultTrackedActionExpiration = 5 * time.Minute
	// DefaultYellowRequestCount is the default tracker yellow request count.
	DefaultYellowRequestCount = 10
	// DefaultYellowRequestPercentage is the default tracker yellow request percentage.
	DefaultYellowRequestPercentage = 0.005 // 0.5% or 50 bps
	// DefaultRedRequestCount is the default tracker red request count.
	DefaultRedRequestCount = 50
	// DefaultRedRequestPercentage is the default tracker yellow request percentage.
	DefaultRedRequestPercentage = 0.05 // 5% or 500 bps
)

// Signal is a status signal.
type Signal string

// Signal constants
const (
	SignalUnknown Signal = ""
	SignalGreen   Signal = "green"
	SignalYellow  Signal = "yellow"
	SignalRed     Signal = "red"
)

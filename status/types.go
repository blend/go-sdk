/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package status

import "time"

//
// API types
//

// FreeformResult is a json result from the freeform check.
type FreeformResult map[string]bool

// TrackedActionsResult are tracked actions details.
type TrackedActionsResult struct {
	Status     Signal          `json:"status"`
	SubSystems map[string]Info `json:"subsystems"`
}

// Info wraps tracked details.
type Info struct {
	Name    string  `json:"name"`
	Status  Signal  `json:"status"`
	Details Details `json:"details"`
}

// Details holds the details about the status results.
type Details struct {
	ErrorCount     int            `json:"errorCount"`
	RequestCount   int            `json:"requestCount"`
	ErrorBreakdown map[string]int `json:"errorBreakdown"`
}

//
// Internal Types
//

// RequestInfo is a type.
type RequestInfo struct {
	RequestTime time.Time
}

// ErrorInfo is a type.
type ErrorInfo struct {
	RequestInfo
	Args interface{}
}

// freeformCheckResult is a result from a freeform check status.
type freeformCheckResult struct {
	ServiceName string
	Ok          bool
	Err         error
}

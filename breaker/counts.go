/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package breaker

// Counts holds the numbers of requests and their successes/failures.
// CircuitBreaker clears the internal Counts either
// on the change of the state or at the closed-state intervals.
// Counts ignores the results of the requests sent before clearing.
type Counts struct {
	Requests		int64	`json:"requests"`
	TotalSuccesses		int64	`json:"totalSuccesses"`
	TotalFailures		int64	`json:"totalFailures"`
	ConsecutiveSuccesses	int64	`json:"consecutiveSuccesses"`
	ConsecutiveFailures	int64	`json:"consecutiveFailures"`
}

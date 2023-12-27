/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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

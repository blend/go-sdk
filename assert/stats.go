/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package assert

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

// assertcount is the total number of assetions run during the package lifetime.
var assertCount int32

// Increment increments the global assertion count.
func Increment() {
	atomic.AddInt32(&assertCount, int32(1))
}

// Count returns the total number of assertions.
func Count() int {
	return int(atomic.LoadInt32(&assertCount))
}

// started is when the package started.
var started time.Time

// Started marks a started time.
func Started() {
	started = time.Now()
}

// Elapsed returns the time since `Started()`
func Elapsed() time.Duration {
	return time.Since(started)
}

// Rate returns the assertions per second.
func Rate() float64 {
	elapsedSeconds := (float64(Elapsed()) / float64(time.Second))
	return float64(atomic.LoadInt32(&assertCount)) / elapsedSeconds
}

// ReportRate writes the rate summary to stdout.
func ReportRate() {
	fmt.Fprintf(os.Stdout, "asserts: %d Δt: %v λ: %0.2f assert/sec\n", Count(), Elapsed(), Rate())
}

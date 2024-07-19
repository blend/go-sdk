/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ratelimiter

import (
	"context"
	"time"
)

// Wait is a type that allows you to throttle actions
// with sleeps based on a desired rate.
type Wait struct {
	NumberOfActions int64
	Quantum         time.Duration
}

// Wait waits for a calculated throttling time based on the input options.
func (w Wait) Wait(ctx context.Context, actions int64, quantum time.Duration) error {
	return w.WaitTimer(ctx, actions, quantum, nil)
}

// WaitTimer waits with a given (re-used) timer reference.
func (w Wait) WaitTimer(ctx context.Context, actions int64, quantum time.Duration, after *time.Timer) error {
	waitFor := w.Calculate(actions, quantum)
	if waitFor < 0 {
		return nil
	}
	if after == nil {
		after = time.NewTimer(waitFor)
	} else {
		after.Reset(waitFor)
	}
	defer after.Stop()
	select {
	case <-ctx.Done():
		return context.Canceled
	case <-after.C:
		return nil
	}
}

// Calculate takes the observed rate and the desired rate, and returns a quantum to sleep for
// that adjusts the observed rate to match the desired rate.
//
// If the observed rate is _lower_ than the desired rate, the returned value will be negative
// and you're free to ignore it.
//
// If the observed rate is _higher_ than the desired rate, a positive duration will be returned
// which you can pass to a `time.Sleep(...)` or similar.
//
// The wait quantum is derrived from the following algebraic steps (where ? is what we're solving for):
//
//	pb/(pq+?) = rb/rq
//	1/(pq+?) = rb/pb*rq
//	pq+? = (pb*rq)/rb
//	? = ((pb*rq)/rb) - pq
func (w Wait) Calculate(actions int64, quantum time.Duration) time.Duration {
	return time.Duration(((actions * int64(w.Quantum)) / w.NumberOfActions) - int64(quantum))
}

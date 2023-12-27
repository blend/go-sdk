/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package retry

import (
	"context"
	"time"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/ex"
)

var (
	_ async.Interceptor = (*Retrier)(nil)
)

// Retry calls an actioner with retries.
func Retry(ctx context.Context, action Actioner, args interface{}, opts ...Option) (interface{}, error) {
	return New(opts...).Intercept(action).Action(ctx, args)
}

// New wraps an actioner with retries.
func New(opts ...Option) *Retrier {
	retrier := Retrier{}
	Defaults(&retrier)
	for _, opt := range opts {
		opt(&retrier)
	}
	return &retrier
}

// Retrier is the retry agent.
type Retrier struct {
	MaxAttempts         uint
	DelayProvider       DelayProvider
	ShouldRetryProvider ShouldRetryProvider
}

// Intercept calls a function and retries on error or if a should retry provider
// is set, based on the should retry result.
func (r Retrier) Intercept(action Actioner) Actioner {
	return ActionerFunc(func(ctx context.Context, args interface{}) (res interface{}, err error) {
		var attempt uint
		var alarm *time.Timer

		// treat maxAttempts == 0 as mostly unbounded
		// we have to keep attempts for the delay provider
		for attempt = 0; (r.MaxAttempts) == 0 || (attempt < r.MaxAttempts); attempt++ {
			func() {
				defer func() {
					if r := recover(); r != nil {
						err = ex.New(r)
					}
				}()
				res, err = action.Action(ctx, args)
			}()
			if err == nil {
				return
			}
			if !r.ShouldRetryProvider(err) {
				return
			}

			// use a (somewhat) persistent alarm reference
			alarm = time.NewTimer(r.DelayProvider(ctx, attempt))
			select {
			case <-ctx.Done():
				alarm.Stop()
				err = context.Canceled
				return
			case <-alarm.C:
				alarm.Stop()
			}
		}
		return
	})
}

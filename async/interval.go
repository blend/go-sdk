/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package async

import (
	"context"
	"time"

	"github.com/blend/go-sdk/ex"
)

/*
NewInterval returns a new worker that runs an action on an interval.

Example:

	iw := NewInterval(func(ctx context.Context) error { return nil }, 500*time.Millisecond)
	go iw.Start()
	<-iw.Started()
*/
func NewInterval(action ContextAction, interval time.Duration, options ...IntervalOption) *Interval {
	i := Interval{
		Latch:		NewLatch(),
		Action:		action,
		Context:	context.Background(),
		Interval:	interval,
	}
	for _, option := range options {
		option(&i)
	}
	return &i
}

// IntervalOption is an option for the interval worker.
type IntervalOption func(*Interval)

// OptIntervalDelay sets the interval worker start delay.
func OptIntervalDelay(d time.Duration) IntervalOption {
	return func(i *Interval) {
		i.Delay = d
	}
}

// OptIntervalContext sets the interval worker context.
func OptIntervalContext(ctx context.Context) IntervalOption {
	return func(i *Interval) {
		i.Context = ctx
	}
}

// OptIntervalStopOnError sets if the interval worker should stop on action error.
func OptIntervalStopOnError(stopOnError bool) IntervalOption {
	return func(i *Interval) {
		i.StopOnError = stopOnError
	}
}

// OptIntervalErrors sets the interval worker start error channel.
func OptIntervalErrors(errors chan error) IntervalOption {
	return func(i *Interval) {
		i.Errors = errors
	}
}

// Interval is a background worker that performs an action on an interval.
type Interval struct {
	*Latch
	Context		context.Context
	Interval	time.Duration
	Action		ContextAction
	Delay		time.Duration
	StopOnError	bool
	Errors		chan error
}

/*
Start starts the worker.

This will start the internal ticker, with a default initial delay of the given interval, and will return an ErrCannotStart if the interval worker is already started.

This call will block.
*/
func (i *Interval) Start() error {
	if !i.CanStart() {
		return ex.New(ErrCannotStart)
	}
	i.Starting()
	return i.Dispatch()
}

// Stop stops the worker.
func (i *Interval) Stop() error {
	if !i.CanStop() {
		return ex.New(ErrCannotStop)
	}
	i.Stopping()
	<-i.NotifyStopped()
	i.Latch.Reset()	// reset the latch in case we have to start again
	return nil
}

// Dispatch is the main dispatch loop.
func (i *Interval) Dispatch() (err error) {
	i.Started()

	if i.Delay > 0 {
		time.Sleep(i.Delay)
	}

	tick := time.NewTicker(i.Interval)
	defer func() {
		tick.Stop()
		i.Stopped()
	}()

	var stopping <-chan struct{}
	for {
		stopping = i.NotifyStopping()
		// check stopping conditions first
		select {
		case <-i.Context.Done():
			return
		case <-stopping:
			return
		default:
		}

		select {
		case <-tick.C:
			err = i.Action(context.Background())
			if err != nil {
				if i.StopOnError {
					return
				}
				if i.Errors != nil {
					i.Errors <- err
				}
			}
		case <-i.Context.Done():
			return
		case <-stopping:
			return
		}
	}
}

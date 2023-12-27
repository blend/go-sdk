/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package async

import (
	"sync"
	"sync/atomic"
)

// NewLatch creates a new latch.
func NewLatch() *Latch {
	l := new(Latch)
	l.Reset()
	return l
}

/*
Latch is a helper to coordinate goroutine lifecycles, specifically waiting for goroutines to start and end.

The lifecycle is generally as follows:

	0 - stopped - goto 1
	1 - starting - goto 2
	2 - started - goto 3
	3 - stopping - goto 0

Control flow is coordinated with chan struct{}, which acts as a semaphore but can only
alert (1) listener as it is buffered.

In order to start a `stopped` latch, you must call `.Reset()` first to initialize channels.
*/
type Latch struct {
	sync.Mutex

	state	int32

	starting	chan struct{}
	started		chan struct{}
	stopping	chan struct{}
	stopped		chan struct{}
}

// Reset resets the latch.
func (l *Latch) Reset() {
	l.Lock()
	atomic.StoreInt32(&l.state, LatchStopped)
	l.starting = make(chan struct{}, 1)
	l.started = make(chan struct{}, 1)
	l.stopping = make(chan struct{}, 1)
	l.stopped = make(chan struct{}, 1)
	l.Unlock()
}

// CanStart returns if the latch can start.
func (l *Latch) CanStart() bool {
	return atomic.LoadInt32(&l.state) == LatchStopped
}

// CanStop returns if the latch can stop.
func (l *Latch) CanStop() bool {
	return atomic.LoadInt32(&l.state) == LatchStarted
}

// IsStarting returns if the latch state is LatchStarting
func (l *Latch) IsStarting() bool {
	return atomic.LoadInt32(&l.state) == LatchStarting
}

// IsStarted returns if the latch state is LatchStarted.
func (l *Latch) IsStarted() bool {
	return atomic.LoadInt32(&l.state) == LatchStarted
}

// IsStopping returns if the latch state is LatchStopping.
func (l *Latch) IsStopping() bool {
	return atomic.LoadInt32(&l.state) == LatchStopping
}

// IsStopped returns if the latch state is LatchStopped.
func (l *Latch) IsStopped() (isStopped bool) {
	return atomic.LoadInt32(&l.state) == LatchStopped
}

// NotifyStarting returns the starting signal.
// It is used to coordinate the transition from stopped -> starting.
// There can only be (1) effective listener at a time for these events.
func (l *Latch) NotifyStarting() (notifyStarting <-chan struct{}) {
	l.Lock()
	notifyStarting = l.starting
	l.Unlock()
	return
}

// NotifyStarted returns the started signal.
// It is used to coordinate the transition from starting -> started.
// There can only be (1) effective listener at a time for these events.
func (l *Latch) NotifyStarted() (notifyStarted <-chan struct{}) {
	l.Lock()
	notifyStarted = l.started
	l.Unlock()
	return
}

// NotifyStopping returns the should stop signal.
// It is used to trigger the transition from running -> stopping -> stopped.
// There can only be (1) effective listener at a time for these events.
func (l *Latch) NotifyStopping() (notifyStopping <-chan struct{}) {
	l.Lock()
	notifyStopping = l.stopping
	l.Unlock()
	return
}

// NotifyStopped returns the stopped signal.
// It is used to coordinate the transition from stopping -> stopped.
// There can only be (1) effective listener at a time for these events.
func (l *Latch) NotifyStopped() (notifyStopped <-chan struct{}) {
	l.Lock()
	notifyStopped = l.stopped
	l.Unlock()
	return
}

// Starting signals the latch is starting.
// This is typically done before you kick off a goroutine.
func (l *Latch) Starting() {
	if l.IsStarting() {
		return
	}
	atomic.StoreInt32(&l.state, LatchStarting)
	l.starting <- struct{}{}
}

// Started signals that the latch is started and has entered
// the `IsStarted` state.
func (l *Latch) Started() {
	if l.IsStarted() {
		return
	}
	atomic.StoreInt32(&l.state, LatchStarted)
	l.started <- struct{}{}
}

// Stopping signals the latch to stop.
// It could also be thought of as `SignalStopping`.
func (l *Latch) Stopping() {
	if l.IsStopping() {
		return
	}
	atomic.StoreInt32(&l.state, LatchStopping)
	l.stopping <- struct{}{}
}

// Stopped signals the latch has stopped.
func (l *Latch) Stopped() {
	if l.IsStopped() {
		return
	}
	atomic.StoreInt32(&l.state, LatchStopped)
	l.stopped <- struct{}{}
}

// WaitStarted triggers `Starting` and waits for the `Started` signal.
func (l *Latch) WaitStarted() {
	if !l.CanStart() {
		return
	}
	started := l.NotifyStarted()
	l.Starting()
	<-started
}

// WaitStopped triggers `Stopping` and waits for the `Stopped` signal.
func (l *Latch) WaitStopped() {
	if !l.CanStop() {
		return
	}
	stopped := l.NotifyStopped()
	l.Stopping()
	<-stopped
}

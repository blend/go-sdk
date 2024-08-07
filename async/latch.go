/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package async

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
alert (1) listener as it is buffered. It is incorrect to use Latch for anything other than
to tie one goroutine to another. Writers *must* perform state transitions in the recommended order,
readers *must* read for state transitions in the expected order. Failure to do so can result in deadlocks.
Channels are used as the synchronization mechanism, we assume a 1-1 tie between goroutines with a single
read/write pair on either end (in any order).

In order to start a `stopped` latch, you must call `.Reset()` first to initialize channels.
*/
type Latch struct {
	state int32

	starting chan struct{}
	started  chan struct{}
	stopping chan struct{}
	stopped  chan struct{}
}

// Reset resets the latch.
func (l *Latch) Reset() {
	l.starting = make(chan struct{}, 1)
	l.started = make(chan struct{}, 1)
	l.stopping = make(chan struct{}, 1)
	l.stopped = make(chan struct{}, 1)
}

// CanStart returns if the latch can start.
func (l *Latch) CanStart() bool {
	return l.state == LatchStopped
}

// CanStop returns if the latch can stop.
func (l *Latch) CanStop() bool {
	return l.state == LatchStarted
}

// IsStarting returns if the latch state is LatchStarting
func (l *Latch) IsStarting() bool {
	return l.state == LatchStarting
}

// IsStarted returns if the latch state is LatchStarted.
func (l *Latch) IsStarted() bool {
	return l.state == LatchStarted
}

// IsStopping returns if the latch state is LatchStopping.
func (l *Latch) IsStopping() bool {
	return l.state == LatchStopping
}

// IsStopped returns if the latch state is LatchStopped.
func (l *Latch) IsStopped() (isStopped bool) {
	return l.state == LatchStopped
}

// NotifyStarting returns the starting signal.
// It is used to coordinate the transition from stopped -> starting.
// There can only be (1) effective listener at a time for these events.
func (l *Latch) NotifyStarting() (notifyStarting <-chan struct{}) {
	notifyStarting = l.starting
	return
}

// NotifyStarted returns the started signal.
// It is used to coordinate the transition from starting -> started.
// There can only be (1) effective listener at a time for these events.
func (l *Latch) NotifyStarted() (notifyStarted <-chan struct{}) {
	notifyStarted = l.started
	return
}

// NotifyStopping returns the should stop signal.
// It is used to trigger the transition from running -> stopping -> stopped.
// There can only be (1) effective listener at a time for these events.
func (l *Latch) NotifyStopping() (notifyStopping <-chan struct{}) {
	notifyStopping = l.stopping
	return
}

// NotifyStopped returns the stopped signal.
// It is used to coordinate the transition from stopping -> stopped.
// There can only be (1) effective listener at a time for these events.
func (l *Latch) NotifyStopped() (notifyStopped <-chan struct{}) {
	notifyStopped = l.stopped
	return
}

// Starting signals the latch is starting.
// This is typically done before you kick off a goroutine.
func (l *Latch) Starting() {
	if l.IsStarting() {
		return
	}
	l.state = LatchStarting
	l.starting <- struct{}{}
}

// Started signals that the latch is started and has entered
// the `IsStarted` state.
func (l *Latch) Started() {
	if l.IsStarted() {
		return
	}
	l.state = LatchStarted
	l.started <- struct{}{}
}

// Stopping signals the latch to stop.
// It could also be thought of as `SignalStopping`.
func (l *Latch) Stopping() {
	if l.IsStopping() {
		return
	}
	l.state = LatchStopping
	l.stopping <- struct{}{}
}

// Stopped signals the latch has stopped.
func (l *Latch) Stopped() {
	if l.IsStopped() {
		return
	}
	l.state = LatchStopped
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

package async

import (
	"sync"
	"sync/atomic"
)

// NewLatch creates a new latch.
func NewLatch() *Latch {
	return &Latch{
		starting: make(chan struct{}),
		resuming: make(chan struct{}),
		started:  make(chan struct{}),
		active:   make(chan struct{}),
		pausing:  make(chan struct{}),
		paused:   make(chan struct{}),
		stopping: make(chan struct{}),
		stopped:  make(chan struct{}),
	}
}

/*
Latch is a helper to coordinate goroutine lifecycles, specifically waiting for goroutines to start and end.

The lifecycle is generally as follows:

	0 - stopped
	1 - starting
	2 - started - goto 3, goto 4, goto 6
	3 - active
	4 - pausing
	5 - paused - goto 6, goto 7
	6 - resuming - goto 2
	7 - stopping - goto 0

Control flow is coordinated with chan struct{}, which allows waiters to pull from the
channel and the triggers to close them.

As a result, each state includes a transition notification, i.e. Starting() triggers <-NotifyStarting().
*/
type Latch struct {
	sync.Mutex
	state int32

	starting chan struct{}
	resuming chan struct{}
	started  chan struct{}
	active   chan struct{}
	pausing  chan struct{}
	paused   chan struct{}
	stopping chan struct{}
	stopped  chan struct{}
}

// CanStart returns if the latch can start.
func (l *Latch) CanStart() bool {
	return atomic.LoadInt32(&l.state) == LatchStopped
}

// CanResume returns if the latch can resume.
func (l *Latch) CanResume() bool {
	return atomic.LoadInt32(&l.state) == LatchPaused
}

// CanPause returns if the latch can pause.
func (l *Latch) CanPause() bool {
	return atomic.LoadInt32(&l.state) == LatchStarted
}

// CanStop returns if the latch can stop.
func (l *Latch) CanStop() bool {
	return atomic.LoadInt32(&l.state) == LatchStarted
}

// IsStarting returns if the latch state is LatchStarting
func (l *Latch) IsStarting() bool {
	return atomic.LoadInt32(&l.state) == LatchStarting
}

// IsResuming returns if the latch state is LatchResuming.
func (l *Latch) IsResuming() bool {
	return atomic.LoadInt32(&l.state) == LatchResuming
}

// IsStarted returns if the latch state is LatchStarted.
func (l *Latch) IsStarted() bool {
	return atomic.LoadInt32(&l.state) == LatchStarted
}

// IsPausing returns if the latch state is LatchPausing.
func (l *Latch) IsPausing() bool {
	return atomic.LoadInt32(&l.state) == LatchPausing
}

// IsPaused returns if the latch state is LatchPaused.
func (l *Latch) IsPaused() bool {
	return atomic.LoadInt32(&l.state) == LatchPaused
}

// IsStopping returns if the latch state is LatchStopping.
func (l *Latch) IsStopping() bool {
	return atomic.LoadInt32(&l.state) == LatchStopping
}

// IsStopped returns if the latch state is LatchStopped.
func (l *Latch) IsStopped() (isStopped bool) {
	return atomic.LoadInt32(&l.state) == LatchStopped
}

// NotifyStarting returns the started signal.
// It is used to coordinate the transition from stopped -> starting.
func (l *Latch) NotifyStarting() (notifyStarting <-chan struct{}) {
	l.Lock()
	notifyStarting = l.starting
	l.Unlock()
	return
}

// NotifyResuming returns the resuming signal.
// It is used to coordinate the transition from paused -> running.
func (l *Latch) NotifyResuming() (notifyResuming <-chan struct{}) {
	l.Lock()
	notifyResuming = l.resuming
	l.Unlock()
	return
}

// NotifyStarted returns the started signal.
// It is used to coordinate the transition from starting -> started.
func (l *Latch) NotifyStarted() (notifyStarted <-chan struct{}) {
	l.Lock()
	notifyStarted = l.started
	l.Unlock()
	return
}

// NotifyActive returns the active signal.
// It is used to coordinate the transition from started -> active.
func (l *Latch) NotifyActive() (notifyActive <-chan struct{}) {
	l.Lock()
	notifyActive = l.active
	l.Unlock()
	return
}

// NotifyPausing returns the pausing signal.
// It is used to coordinate the transition from running -> pausing.
func (l *Latch) NotifyPausing() (notifyPausing <-chan struct{}) {
	l.Lock()
	notifyPausing = l.pausing
	l.Unlock()
	return
}

// NotifyPaused returns the paused signal.
// It is used to coordinate the transition from pausing -> paused.
func (l *Latch) NotifyPaused() (notifyPaused <-chan struct{}) {
	l.Lock()
	notifyPaused = l.paused
	l.Unlock()
	return
}

// NotifyStopping returns the should stop signal.
// It is used to trigger the transition from running -> stopping -> stopped.
func (l *Latch) NotifyStopping() (notifyStopping <-chan struct{}) {
	l.Lock()
	notifyStopping = l.stopping
	l.Unlock()
	return
}

// NotifyStopped returns the stopped signal.
// It is used to coordinate the transition from stopping -> stopped.
func (l *Latch) NotifyStopped() (notifyStopped <-chan struct{}) {
	l.Lock()
	notifyStopped = l.stopped
	l.Unlock()
	return
}

// Starting signals the latch is starting.
// This is typically done before you kick off a goroutine.
func (l *Latch) Starting() {
	l.Lock()
	defer l.Unlock()

	if l.IsStarting() {
		return
	}
	atomic.StoreInt32(&l.state, LatchStarting)
	close(l.starting)
	l.starting = make(chan struct{})
}

// Resuming signals that the latch is resuming and has entered
// the `IsResuming` state.
func (l *Latch) Resuming() {
	l.Lock()
	defer l.Unlock()

	if l.IsResuming() {
		return
	}
	atomic.StoreInt32(&l.state, LatchResuming)
	close(l.resuming)
	l.resuming = make(chan struct{})
}

// Started signals that the latch is started and has entered
// the `IsStarted` state.
func (l *Latch) Started() {
	l.Lock()
	defer l.Unlock()

	if l.IsStarted() {
		return
	}
	atomic.StoreInt32(&l.state, LatchStarted)
	close(l.started)
	l.started = make(chan struct{})
}

// Pausing signals that the latch is pausing and has entered
// the `IsPausing` state.
func (l *Latch) Pausing() {
	l.Lock()
	defer l.Unlock()

	if l.IsPausing() {
		return
	}
	atomic.StoreInt32(&l.state, LatchPausing)
	close(l.pausing)
	l.pausing = make(chan struct{})
}

// Paused signals that the latch is paused and has entered
// the `IsPaused` state.
func (l *Latch) Paused() {
	l.Lock()
	defer l.Unlock()

	if l.IsPaused() {
		return
	}
	atomic.StoreInt32(&l.state, LatchPaused)
	close(l.paused)
	l.paused = make(chan struct{})
}

// Stopping signals the latch to stop.
// It could also be thought of as `SignalStopping`.
func (l *Latch) Stopping() {
	l.Lock()
	defer l.Unlock()

	if l.IsStopping() {
		return
	}
	atomic.StoreInt32(&l.state, LatchStopping)
	close(l.stopping)
	l.stopping = make(chan struct{})
}

// Stopped signals the latch has stopped.
func (l *Latch) Stopped() {
	l.Lock()
	defer l.Unlock()

	if l.IsStopped() {
		return
	}
	atomic.StoreInt32(&l.state, LatchStopped)
	close(l.stopped)
	l.stopped = make(chan struct{})
}

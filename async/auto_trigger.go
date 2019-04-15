package async

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/blend/go-sdk/ex"
)

/*
NewAutoTrigger returns a new singleton that can be used to automatically trigger an action.

The action will be called after a given quantum or a given number of increments have happened.

This can be useful to debounce calls to upstream systems, such as a checkpointer.
*/
func NewAutoTrigger(action ContextAction, interval time.Duration, maxCount int, options ...AutoTriggerOption) *AutoTrigger {
	at := AutoTrigger{
		Latch:         NewLatch(),
		Action:        action,
		MaxCount:      int32(maxCount),
		Interval:      interval,
		Context:       context.Background(),
		TriggerOnStop: true,
	}
	for _, option := range options {
		option(&at)
	}
	return &at
}

// AutoTriggerOption is an option for an auto-action.
type AutoTriggerOption func(*AutoTrigger)

// OptAutoTriggerMaxCount sets the auto-action max count.
func OptAutoTriggerMaxCount(maxCount int32) AutoTriggerOption {
	return func(aa *AutoTrigger) {
		aa.MaxCount = maxCount
	}
}

// OptAutoTriggerInterval sets the auto-action interval.
func OptAutoTriggerInterval(d time.Duration) AutoTriggerOption {
	return func(aa *AutoTrigger) {
		aa.Interval = d
	}
}

// OptAutoTriggerErrors sets the auto-action error channel.
func OptAutoTriggerErrors(errors chan error) AutoTriggerOption {
	return func(aa *AutoTrigger) {
		aa.Errors = errors
	}
}

// OptAutoTriggerOnStop sets if the auto-action should call the action on shutdown.
func OptAutoTriggerOnStop(errors chan error) AutoTriggerOption {
	return func(aa *AutoTrigger) {
		aa.Errors = errors
	}
}

// AutoTrigger triggers an action on a given interval or after a given number of increments.
type AutoTrigger struct {
	*Latch

	Action        ContextAction
	Context       context.Context
	Errors        chan error
	Interval      time.Duration
	MaxCount      int32
	TriggerOnStop bool

	Counter int32
}

// Background returns a background context.
func (a *AutoTrigger) Background() context.Context {
	if a.Context != nil {
		return a.Context
	}
	return context.Background()
}

/*
Start starts the singleton.

This call blocks. To call it asynchronously:

	go a.Start()
	<-a.NotifyStarted()

This will start the singleton and wait for it to enter the running state.
*/
func (a *AutoTrigger) Start() error {
	if !a.CanStart() {
		return ex.New(ErrCannotStart)
	}
	a.Starting()
	a.Dispatch()
	return nil
}

// Dispatch is the main run loop.
func (a *AutoTrigger) Dispatch() {
	a.Started()
	ticker := time.Tick(a.Interval)
	for {
		select {
		case <-ticker:
			a.Trigger(a.Background())
		case <-a.NotifyPausing():
			a.Paused()
			select {
			case <-a.NotifyResuming():
				a.Started()
			case <-a.NotifyStopping():
				a.Stopped()
				return
			}
		case <-a.NotifyStopping():
			if a.TriggerOnStop {
				a.Trigger(a.Background())
			}
			a.Stopped()
			return
		}
	}
}

// Stop stops the auto-action singleton.
func (a *AutoTrigger) Stop() error {
	if !a.CanStop() {
		return ex.New(ErrCannotStop)
	}
	a.Stopping()
	<-a.NotifyStopped()
	return nil
}

// Increment updates the count
func (a *AutoTrigger) Increment(ctx context.Context) {
	if atomic.CompareAndSwapInt32(&a.Counter, a.MaxCount-1, 0) {
		a.Trigger(ctx)
		return
	}
	atomic.AddInt32(&a.Counter, 1)
}

// Trigger invokes the action if one is set, it will acquire the lock and hold it for the duration of the call to the action.
func (a *AutoTrigger) Trigger(ctx context.Context) {
	a.Lock()
	defer a.Unlock()
	defer func() {
		if r := recover(); r != nil {
			if a.Errors != nil {
				a.Errors <- ex.New(r)
			}
		}
	}()

	if err := a.Action(ctx); err != nil && a.Errors != nil {
		a.Errors <- err
	}
}

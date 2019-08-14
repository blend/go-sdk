package logger

import (
	"context"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/ex"
)

// NewWorker returns a new worker.
func NewWorker(listener Listener) *Worker {
	return &Worker{
		Latch:    async.NewLatch(),
		Listener: listener,
		Work:     make(chan EventWithContext, DefaultWorkerQueueDepth),
	}
}

// Worker is an agent that processes a listener.
type Worker struct {
	*async.Latch
	Errors   chan error
	Listener Listener
	Work     chan EventWithContext
}

// Start starts the worker.
func (w *Worker) Start() error {
	if !w.CanStart() {
		return ex.New(async.ErrCannotStart)
	}
	w.Starting()
	w.Dispatch()
	return nil
}

// Dispatch is the main listen loop
func (w *Worker) Dispatch() {
	w.Started()

	var e EventWithContext
	var err error
	var pausing <-chan struct{}
	var stopping <-chan struct{}

	for {
		pausing = w.NotifyPausing()
		stopping = w.NotifyStopping()

		// we have to do this effectively twice to add precedence to the stop and pause
		// signals
		select {
		case <-pausing:
			w.Paused()
			select {
			case <-w.NotifyResuming():
				w.Reset()
				w.Started()
				continue
			case <-w.NotifyStopping():
				w.Stopped()
				return
			}
		case <-stopping:
			w.Stopped()
			return
		default:
		}

		select {
		case <-pausing:
			w.Paused()
			select {
			case <-w.NotifyResuming():
				w.Reset()
				w.Started()
				continue
			case <-w.NotifyStopping():
				w.Stopped()
				return
			}
		case <-stopping:
			w.Stopped()
			return
		case e = <-w.Work:
			if err = w.Process(e); err != nil && w.Errors != nil {
				w.Errors <- err
			}
		}
	}
}

// Process calls the listener for an event.
func (w *Worker) Process(ec EventWithContext) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ex.New(r)
			return
		}
	}()
	w.Listener(ec.Context, ec.Event)
	return
}

// DrainContext pauses the worker and synchronously processes any remaining work.
// It then restarts the worker.
func (w *Worker) DrainContext(ctx context.Context) error {
	// if the worker is currently processing work, wait for it to finish.
	notifyPaused := w.NotifyPaused()
	w.Pausing()
	select {
	case <-notifyPaused:
		break
	case <-ctx.Done():
		return context.Canceled
	}
	defer func() {
		w.Resuming()
	}()

	var work EventWithContext
	var err error

	workLeft := len(w.Work)
	done := make(chan struct{})
	go func() {
		defer close(done)
		for index := 0; index < workLeft; index++ {
			work = <-w.Work
			if err = w.Process(work); err != nil && w.Errors != nil {
				w.Errors <- err
			}
		}
	}()

	select {
	case <-ctx.Done():
		return context.Canceled
	case <-done:
		return nil
	}
}

// Stop stops the worker.
func (w *Worker) Stop() error {
	if !w.CanStop() {
		return ex.New(async.ErrCannotStop)
	}
	if w.IsActive() {
		<-w.NotifyStarted()
	}

	w.Stopping()
	<-w.NotifyStopped()

	var work EventWithContext
	var err error

	workLeft := len(w.Work)
	for index := 0; index < workLeft; index++ {
		work = <-w.Work
		if err = w.Process(work); err != nil && w.Errors != nil {
			w.Errors <- err
		}
	}
	return nil
}

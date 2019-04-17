package logger

import (
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
	for {
		select {
		case e = <-w.Work:
			if err = w.Process(e); err != nil && w.Errors != nil {
				w.Errors <- err
			}
		case <-w.NotifyPausing():
			w.Paused()
			select {
			case <-w.NotifyResuming():
				w.Started()
			case <-w.NotifyStopping():
				w.Stopped()
				return
			}
		case <-w.NotifyStopping():
			w.Stopped()
			return
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

// Drain stops the worker and synchronously processes any remaining work.
// It then restarts the worker.
func (w *Worker) Drain() {
	if w.CanPause() {
		w.Pausing()
		defer func() {
			w.Resuming()
		}()
	}

	var work EventWithContext
	var err error
	workLeft := len(w.Work)
	for index := 0; index < workLeft; index++ {
		work = <-w.Work
		if err = w.Process(work); err != nil && w.Errors != nil {
			w.Errors <- err
		}
	}
}

// Stop stops the worker.
func (w *Worker) Stop() error {
	if !w.CanStop() {
		return ex.New(async.ErrCannotStop)
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

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
	w.DispatchWork()
	return nil
}

// Stop stops the worker.
func (w *Worker) Stop() error {
	return w.StopContext(context.Background())
}

// StopContext stops the worker and processe what work is left in the queue.
// The worker will be stopped at the end, and it will be required to Start
// the worker again to
func (w *Worker) StopContext(ctx context.Context) error {
	// if the worker is currently processing work, wait for it to finish.
	notifyStopped := w.NotifyStopped()

	// signal that the DispatchWork loop should stop.
	w.Stopping()

	// wait for the DispatchWork loop to stop
	// but also check if we've timed out.
	select {
	case <-notifyStopped:
		break
	case <-ctx.Done():
		return context.Canceled
	}

	// process what's left of the work queue.
	var work EventWithContext
	var err error

	workLeft := len(w.Work)
	drained := make(chan struct{})
	go func() {
		// notify once the last of the work
		// in the queue is complete.
		defer close(drained)

		// go through what's left of the work
		// be mindful of the outer context timeout.
		for index := 0; index < workLeft; index++ {
			select {
			case <-ctx.Done():
				w.Errors <- context.Canceled
				return
			case work = <-w.Work:
				if err = w.Process(work); err != nil && w.Errors != nil {
					w.Errors <- err
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		return context.Canceled
	case <-drained:
		return nil
	}
}

// AbortContext stops the worker but does not process the remaining work.
func (w *Worker) AbortContext(ctx context.Context) error {
	notifyStopped := w.NotifyStopped()
	w.Stopping()
	select {
	case <-notifyStopped:
		return nil
	case <-ctx.Done():
		return context.Canceled
	}
}

// DispatchWork is the dispatch loop where
// work is processed.
func (w *Worker) DispatchWork() {
	w.Started()
	var e EventWithContext
	var err error

	notifyStopping := w.NotifyStopping()
	for {
		select {
		case <-notifyStopping:
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

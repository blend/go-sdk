package async

import (
	"context"

	"github.com/blend/go-sdk/ex"
)

// NewWorker creates a new worker.
func NewWorker(action WorkAction) *Worker {
	return &Worker{
		Context: context.Background(),
		Latch:   NewLatch(),
		Action:  action,
		Work:    make(chan interface{}),
	}
}

// Worker is a worker that is pushed work over a channel.
// It is used by other work distribution types (i.e. queue and batch)
// but can also be used independently.
type Worker struct {
	Latch *Latch

	Context   context.Context
	Action    WorkAction
	Finalizer WorkerFinalizer

	Errors chan error
	Work   chan interface{}
}

// Background returns the queue worker background context.
func (w *Worker) Background() context.Context {
	if w.Context != nil {
		return w.Context
	}
	return context.Background()
}

// NotifyStarted returns the underlying latch signal.
func (w *Worker) NotifyStarted() <-chan struct{} {
	return w.Latch.NotifyStarted()
}

// NotifyStopped returns the underlying latch signal.
func (w *Worker) NotifyStopped() <-chan struct{} {
	return w.Latch.NotifyStarted()
}

// Enqueue adds an item to the work queue.
func (w *Worker) Enqueue(obj interface{}) {
	w.Work <- obj
}

// Start starts the worker with a given context.
func (w *Worker) Start() error {
	if !w.Latch.CanStart() {
		return ex.New(ErrCannotStart)
	}
	w.Latch.Starting()
	w.Dispatch()
	return nil
}

// Dispatch starts the listen loop for work.
func (w *Worker) Dispatch() {
	w.Latch.Started()
	var workItem interface{}
	var stopping <-chan struct{}
	for {
		stopping = w.Latch.NotifyStopping()
		// we should always check stopped
		// before also blocking on work or stopping
		select {
		case <-stopping:
			w.Latch.Stopped()
			return
		case <-w.Background().Done():
			w.Latch.Stopped()
			return
		default:
		}

		// block on work or stopping
		select {
		case workItem = <-w.Work:
			w.Execute(w.Background(), workItem)
		case <-stopping:
			w.Latch.Stopped()
			return
		case <-w.Background().Done():
			w.Latch.Stopped()
			return
		}
	}
}

// Execute invokes the action and recovers panics.
func (w *Worker) Execute(ctx context.Context, workItem interface{}) {
	defer func() {
		if r := recover(); r != nil {
			w.HandleError(ex.New(r))
		}
		if w.Finalizer != nil {
			w.HandleError(w.Finalizer(ctx, w))
		}
	}()
	if w.Action != nil {
		w.HandleError(w.Action(ctx, workItem))
	}
}

// Stop stop the worker.
// The work left in the queue will remain.
func (w *Worker) Stop() error {
	if !w.Latch.CanStop() {
		return ex.New(ErrCannotStop)
	}
	w.Latch.Stopping()
	<-w.Latch.NotifyStopped()
	return nil
}

// Drain stops the worker and synchronously waits
// for in progress items to finish.
func (w *Worker) Drain(ctx context.Context) {
	drainComplete := make(chan struct{})
	go func() {
		defer close(drainComplete)
		w.Latch.Stopping()
		<-w.Latch.NotifyStopped()
	}()

	select {
	case <-drainComplete:
		return
	case <-ctx.Done():
		return
	}
}

// HandleError sends a non-nil err to the error
// collector if one is provided.
func (w *Worker) HandleError(err error) {
	if err == nil {
		return
	}
	if w.Errors == nil {
		return
	}
	w.Errors <- err
}

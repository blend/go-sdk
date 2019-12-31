package async

import (
	"context"

	"github.com/blend/go-sdk/ex"
)

// NewWorker creates a new worker.
func NewWorker(action WorkAction) *Worker {
	return &Worker{
		Latch:  NewLatch(),
		Action: action,
		Work:   make(chan interface{}),
	}
}

// Worker is a worker that is pushed work over a channel.
// It is used by other work distribution types (i.e. queue and batch)
// but can also be used independently.
type Worker struct {
	Latch     *Latch
	Context   context.Context
	Action    WorkAction
	Finalizer WorkerFinalizer
	Errors    chan error
	Work      chan interface{}
}

// Background returns the queue worker background context.
func (qw *Worker) Background() context.Context {
	if qw.Context != nil {
		return qw.Context
	}
	return context.Background()
}

// NotifyStarted returns the underlying latch signal.
func (qw *Worker) NotifyStarted() <-chan struct{} {
	return qw.Latch.NotifyStarted()
}

// NotifyStopped returns the underlying latch signal.
func (qw *Worker) NotifyStopped() <-chan struct{} {
	return qw.Latch.NotifyStarted()
}

// Enqueue adds an item to the work queue.
func (qw *Worker) Enqueue(obj interface{}) {
	qw.Work <- obj
}

// Start starts the worker with a given context.
func (qw *Worker) Start() error {
	if !qw.Latch.CanStart() {
		return ex.New(ErrCannotStart)
	}
	qw.Latch.Starting()
	qw.Dispatch()
	return nil
}

// Dispatch starts the listen loop for work.
func (qw *Worker) Dispatch() {
	qw.Latch.Started()
	var workItem interface{}
	var stopping <-chan struct{}
	for {
		stopping = qw.Latch.NotifyStopping()
		select {
		case workItem = <-qw.Work:
			qw.Execute(qw.Background(), workItem)
		case <-stopping:
			qw.Latch.Stopped()
			return
		}
	}
}

// Execute invokes the action and recovers panics.
func (qw *Worker) Execute(ctx context.Context, workItem interface{}) {
	defer func() {
		if r := recover(); r != nil {
			qw.HandleError(ex.New(r))
		}
		if qw.Finalizer != nil {
			qw.HandleError(qw.Finalizer(ctx, qw))
		}
	}()
	if qw.Action != nil {
		qw.HandleError(qw.Action(ctx, workItem))
	}
}

// Stop stop the worker.
// The work left in the queue will remain.
func (qw *Worker) Stop() error {
	if !qw.Latch.CanStop() {
		return ex.New(ErrCannotStop)
	}
	qw.Latch.Stopping()
	<-qw.Latch.NotifyStopped()
	return nil
}

// Drain stops the worker and synchronously drains the the remaining work
// with a given context.
func (qw *Worker) Drain(ctx context.Context) {
	qw.Latch.Stopping()
	<-qw.Latch.NotifyStopped()

	// create a signal that we've completed draining.
	stopped := make(chan struct{})
	remaining := len(qw.Work)
	go func() {
		defer close(stopped)
		for x := 0; x < remaining; x++ {
			qw.Execute(qw.Background(), <-qw.Work)
		}
	}()
	<-stopped
}

// HandleError sends a non-nil err to the error
// collector if one is provided.
func (qw *Worker) HandleError(err error) {
	if err == nil {
		return
	}
	if qw.Errors == nil {
		return
	}
	qw.Errors <- err
}

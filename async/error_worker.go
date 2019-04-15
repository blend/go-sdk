package async

import (
	"context"

	"github.com/blend/go-sdk/ex"
)

// NewErrorWorker creates a new worker.
func NewErrorWorker(action ErrorWorkAction) *ErrorWorker {
	return &ErrorWorker{
		Latch:  NewLatch(),
		Action: action,
		Work:   make(chan error),
	}
}

// ErrorWorker is a worker that is pushed work over a channel.
type ErrorWorker struct {
	*Latch
	Context   context.Context
	Action    ErrorWorkAction
	Finalizer ErrorWorkerFinalizer
	Work      chan error
}

// Background returns the queue worker background context.
func (qw *ErrorWorker) Background() context.Context {
	if qw.Context != nil {
		return qw.Context
	}
	return context.Background()
}

// Enqueue adds an item to the work queue.
func (qw *ErrorWorker) Enqueue(obj error) {
	qw.Work <- obj
}

// Start starts the worker with a given context.
func (qw *ErrorWorker) Start() error {
	if !qw.CanStart() {
		return ex.New(ErrCannotStart)
	}
	qw.Starting()
	qw.Dispatch()
	return nil
}

// Dispatch starts the listen loop for work.
func (qw *ErrorWorker) Dispatch() {
	qw.Started()
	var workItem error
	for {
		select {
		case workItem = <-qw.Work:
			qw.Execute(qw.Background(), workItem)
		case <-qw.NotifyPausing():
			qw.Paused()
			select {
			case <-qw.NotifyResuming():
				qw.Started()
			case <-qw.NotifyStopping():
				qw.Stopped()
				return
			}
		case <-qw.NotifyStopping():
			qw.Stopped()
			return
		}
	}
}

// Execute invokes the action and recovers panics.
func (qw *ErrorWorker) Execute(ctx context.Context, workItem error) {
	defer func() {
		if qw.Finalizer != nil {
			qw.Finalizer(ctx, qw)
		}
	}()
	if qw.Action != nil {
		qw.Action(ctx, workItem)
	}

}

// Stop stop the worker.
// The work left in the queue will remain.
func (qw *ErrorWorker) Stop() error {
	if !qw.CanStop() {
		return ex.New(ErrCannotStop)
	}
	qw.Stopping()
	<-qw.NotifyStopped()
	return nil
}

// Drain stops the worker and synchronously drains the the remaining work
// with a given context.
func (qw *ErrorWorker) Drain(ctx context.Context) {
	qw.Stopping()
	<-qw.NotifyStopped()

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

// Close stops the worker and cleans up resources.
func (qw *ErrorWorker) Close() error {
	qw.Stopping()
	<-qw.NotifyStopped()
	close(qw.Work)
	return nil
}

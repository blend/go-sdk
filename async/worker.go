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
type Worker struct {
	*Latch
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

// Enqueue adds an item to the work queue.
func (qw *Worker) Enqueue(obj interface{}) {
	qw.Work <- obj
}

// Start starts the worker with a given context.
func (qw *Worker) Start() error {
	if !qw.CanStart() {
		return ex.New(ErrCannotStart)
	}
	qw.Starting()
	qw.Dispatch()
	return nil
}

// Dispatch starts the listen loop for work.
func (qw *Worker) Dispatch() {
	qw.Started()

	var workItem interface{}
	var pausing <-chan struct{}
	var stopping <-chan struct{}

	for {
		pausing = qw.NotifyPausing()
		stopping = qw.NotifyStopping()

		select {
		case workItem = <-qw.Work:
			qw.Execute(qw.Background(), workItem)
		case <-pausing:
			qw.Paused()
			select {
			case <-qw.NotifyResuming():
				qw.Started()
			case <-qw.NotifyStopping():
				qw.Stopped()
				return
			}
		case <-stopping:
			qw.Stopped()
			return
		}
	}
}

// Execute invokes the action and recovers panics.
func (qw *Worker) Execute(ctx context.Context, workItem interface{}) {
	defer func() {
		if r := recover(); r != nil {
			if qw.Errors != nil {
				qw.Errors <- ex.New(r)
			}
		}
		if qw.Finalizer != nil {
			if err := qw.Finalizer(ctx, qw); err != nil {
				if qw.Errors != nil {
					qw.Errors <- ex.New(err)
				}
			}
		}
	}()
	if qw.Action != nil {
		if err := qw.Action(ctx, workItem); err != nil {
			if qw.Errors != nil {
				qw.Errors <- ex.New(err)
			}
		}
	}

}

// Stop stop the worker.
// The work left in the queue will remain.
func (qw *Worker) Stop() error {
	if !qw.CanStop() {
		return ex.New(ErrCannotStop)
	}
	qw.Stopping()
	<-qw.NotifyStopped()
	return nil
}

// Drain stops the worker and synchronously drains the the remaining work
// with a given context.
func (qw *Worker) Drain(ctx context.Context) {
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
func (qw *Worker) Close() error {
	qw.Stopping()
	<-qw.NotifyStopped()
	close(qw.Work)
	return nil
}

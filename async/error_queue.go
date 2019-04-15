package async

import (
	"context"
	"runtime"

	"github.com/blend/go-sdk/ex"
)

// NewErrorQueue returns a new parallel error queue worker.
// It is meant to provide a way to get errors out of other parallel workers and process them.
func NewErrorQueue(action ErrorWorkAction, options ...ErrorQueueOption) *ErrorQueue {
	q := &ErrorQueue{
		Latch:       NewLatch(),
		Action:      action,
		Context:     context.Background(),
		MaxWork:     DefaultQueueMaxWork,
		Parallelism: runtime.NumCPU(),
	}
	for _, option := range options {
		option(q)
	}
	return q
}

// ErrorQueueOption is an option for the error queue worker.
type ErrorQueueOption func(*ErrorQueue)

// OptErrorQueueMaxWork sets the queue worker count.
func OptErrorQueueMaxWork(maxWork int) ErrorQueueOption {
	return func(q *ErrorQueue) {
		q.MaxWork = maxWork
	}
}

// OptErrorQueueParallelism sets the queue worker count.
func OptErrorQueueParallelism(parallelism int) ErrorQueueOption {
	return func(q *ErrorQueue) {
		q.Parallelism = parallelism
	}
}

// OptErrorQueueContext sets the queue worker context.
func OptErrorQueueContext(ctx context.Context) ErrorQueueOption {
	return func(q *ErrorQueue) {
		q.Context = ctx
	}
}

// ErrorQueue is an error queue with multiple workers..
type ErrorQueue struct {
	*Latch

	Action      ErrorWorkAction
	Context     context.Context
	MaxWork     int
	Parallelism int

	// these will typically be set by Start
	Workers chan *ErrorWorker
	Work    chan error
}

// Background returns a background context.
func (pq *ErrorQueue) Background() context.Context {
	if pq.Context != nil {
		return pq.Context
	}
	return context.Background()
}

// Enqueue adds an item to the work queue.
func (pq *ErrorQueue) Enqueue(obj error) {
	pq.Work <- obj
}

// Start starts the queue and its workers.
// This call blocks.
func (pq *ErrorQueue) Start() error {
	if !pq.CanStart() {
		return ex.New(ErrCannotStart)
	}
	pq.Starting()

	// create channel(s)
	if pq.Work == nil {
		pq.Work = make(chan error, pq.MaxWork)
	}
	if pq.Workers == nil {
		pq.Workers = make(chan *ErrorWorker, pq.Parallelism)
	}
	for x := 0; x < pq.Parallelism; x++ {
		worker := NewErrorWorker(pq.Action)
		worker.Context = pq.Context
		worker.Finalizer = pq.ReturnWorker
		go worker.Start()
		<-worker.NotifyStarted()
		pq.Workers <- worker
	}
	pq.Dispatch()
	return nil
}

// Dispatch processes work items in a loop.
func (pq *ErrorQueue) Dispatch() {
	pq.Started()
	var workItem error
	var worker *ErrorWorker
	for {
		select {
		case workItem = <-pq.Work:
			select {
			case worker = <-pq.Workers:
				worker.Enqueue(workItem)
			case <-pq.NotifyStopping():
				pq.Stopped()
				return
			}
		case <-pq.NotifyPausing():
			pq.Paused()
			select {
			case <-pq.NotifyResuming():
				pq.Started()
			case <-pq.NotifyStopping():
				pq.Stopped()
				return
			}
		case <-pq.NotifyStopping():
			pq.Stopped()
			return
		}
	}
}

// Stop stops the queue
func (pq *ErrorQueue) Stop() error {
	if !pq.CanStop() {
		return ex.New(ErrCannotStop)
	}
	for x := 0; x < pq.Parallelism; x++ {
		worker := <-pq.Workers
		worker.Stop()
		pq.Workers <- worker
	}
	return nil
}

// Close stops the queue.
// Any work left in the queue will be discarded.
func (pq *ErrorQueue) Close() error {
	pq.Stopping()
	<-pq.NotifyStopped()
	return nil
}

// ReturnWorker creates an action handler that returns a given worker to the worker queue.
// It wraps any action provided to the queue.
func (pq *ErrorQueue) ReturnWorker(ctx context.Context, worker *ErrorWorker) {
	pq.Workers <- worker
}

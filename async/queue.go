package async

import (
	"context"
	"runtime"

	"github.com/blend/go-sdk/ex"
)

// NewQueue returns a new parallel queue.
func NewQueue(action WorkAction, options ...QueueOption) *Queue {
	q := Queue{
		Latch:       NewLatch(),
		Action:      action,
		Context:     context.Background(),
		MaxWork:     DefaultQueueMaxWork,
		Parallelism: runtime.NumCPU(),
	}
	for _, option := range options {
		option(&q)
	}
	return &q
}

// QueueOption is an option for the queue worker.
type QueueOption func(*Queue)

// OptQueueParallelism sets the queue worker parallelism.
func OptQueueParallelism(parallelism int) QueueOption {
	return func(q *Queue) {
		q.Parallelism = parallelism
	}
}

// OptQueueMaxWork sets the queue worker max work.
func OptQueueMaxWork(maxWork int) QueueOption {
	return func(q *Queue) {
		q.MaxWork = maxWork
	}
}

// OptQueueErrors sets the queue worker start error channel.
func OptQueueErrors(errors chan error) QueueOption {
	return func(q *Queue) {
		q.Errors = errors
	}
}

// OptQueueContext sets the queue worker context.
func OptQueueContext(ctx context.Context) QueueOption {
	return func(q *Queue) {
		q.Context = ctx
	}
}

// Queue is a queue with multiple workers.
type Queue struct {
	*Latch

	Action      WorkAction
	Context     context.Context
	Errors      chan error
	Parallelism int
	MaxWork     int

	// these will typically be set by Start
	Workers chan *Worker
	Work    chan interface{}
}

// Background returns a background context.
func (pq *Queue) Background() context.Context {
	if pq.Context != nil {
		return pq.Context
	}
	return context.Background()
}

// Enqueue adds an item to the work queue.
func (pq *Queue) Enqueue(obj interface{}) {
	pq.Work <- obj
}

// Start starts the queue and its workers.
// This call blocks.
func (pq *Queue) Start() error {
	if !pq.CanStart() {
		return ex.New(ErrCannotStart)
	}
	pq.Starting()

	// create channel(s)
	pq.Work = make(chan interface{}, pq.MaxWork)
	pq.Workers = make(chan *Worker, pq.Parallelism)

	for x := 0; x < pq.Parallelism; x++ {
		worker := NewWorker(pq.Action)
		worker.Context = pq.Context
		worker.Errors = pq.Errors
		worker.Finalizer = pq.ReturnWorker

		// start the worker on its own goroutine
		go worker.Start()
		<-worker.NotifyStarted()
		pq.Workers <- worker
	}
	pq.Dispatch()
	return nil
}

// Dispatch processes work items in a loop.
func (pq *Queue) Dispatch() {
	pq.Started()
	var workItem interface{}
	var worker *Worker
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
func (pq *Queue) Stop() error {
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
func (pq *Queue) Close() error {
	pq.Stopping()
	<-pq.NotifyStopped()
	return nil
}

// ReturnWorker creates an action handler that returns a given worker to the worker queue.
// It wraps any action provided to the queue.
func (pq *Queue) ReturnWorker(ctx context.Context, worker *Worker) error {
	pq.Workers <- worker
	return nil
}

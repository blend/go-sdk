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
func (q *Queue) Background() context.Context {
	if q.Context != nil {
		return q.Context
	}
	return context.Background()
}

// Enqueue adds an item to the work queue.
func (q *Queue) Enqueue(obj interface{}) {
	q.Work <- obj
}

// Start starts the queue and its workers.
// This call blocks.
func (q *Queue) Start() error {
	if !q.CanStart() {
		return ex.New(ErrCannotStart)
	}
	q.Starting()

	// create channel(s)
	q.Work = make(chan interface{}, q.MaxWork)
	q.Workers = make(chan *Worker, q.Parallelism)

	for x := 0; x < q.Parallelism; x++ {
		worker := NewWorker(q.Action)
		worker.Context = q.Context
		worker.Errors = q.Errors
		worker.Finalizer = q.ReturnWorker

		// start the worker on its own goroutine
		go worker.Start()
		<-worker.NotifyStarted()
		q.Workers <- worker
	}
	q.Dispatch()
	return nil
}

// Dispatch processes work items in a loop.
func (q *Queue) Dispatch() {
	q.Started()
	var workItem interface{}
	var worker *Worker
	var stopping <-chan struct{}
	for {
		stopping = q.NotifyStopping()
		select {
		case workItem = <-q.Work:
			stopping = q.NotifyStopping()
			select {
			case worker = <-q.Workers:
				worker.Enqueue(workItem)
			case <-stopping:
				q.Stopped()
				return
			}
		case <-stopping:
			q.Stopped()
			return
		}
	}
}

// Stop stops the queue
func (q *Queue) Stop() error {
	if !q.CanStop() {
		return ex.New(ErrCannotStop)
	}
	q.WaitStopped()
	for x := 0; x < q.Parallelism; x++ {
		worker := <-q.Workers
		worker.Stop()
		q.Workers <- worker
	}
	return nil
}

// Close stops the queue.
// Any work left in the queue will be discarded.
func (q *Queue) Close() error {
	q.WaitStopped()
	return nil
}

// ReturnWorker creates an action handler that returns a given worker to the worker queue.
// It wraps any action provided to the queue.
func (q *Queue) ReturnWorker(ctx context.Context, worker *Worker) error {
	q.Workers <- worker
	return nil
}

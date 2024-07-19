/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package async

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/blend/go-sdk/ex"
)

// NewQueue returns a new parallel queue.
func NewQueue(action WorkAction, options ...QueueOption) *Queue {
	q := Queue{
		Latch:               NewLatch(),
		Action:              action,
		Context:             context.Background(),
		MaxWork:             DefaultQueueMaxWork,
		Parallelism:         runtime.NumCPU(),
		ShutdownGracePeriod: DefaultShutdownGracePeriod,
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

	Action              WorkAction
	Context             context.Context
	Errors              chan error
	Parallelism         int
	MaxWork             int
	ShutdownGracePeriod time.Duration

	// these will typically be set by Start
	AvailableWorkers chan *Worker
	Workers          []*Worker
	Work             chan interface{}
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
	if !q.Latch.CanStart() {
		return ex.New(ErrCannotStart)
	}
	q.Latch.Starting()

	// create channel(s)
	q.Work = make(chan interface{}, q.MaxWork)
	q.AvailableWorkers = make(chan *Worker, q.Parallelism)
	q.Workers = make([]*Worker, q.Parallelism)

	for x := 0; x < q.Parallelism; x++ {
		worker := NewWorker(q.Action)
		worker.Context = q.Context
		worker.Errors = q.Errors
		worker.Finalizer = q.ReturnWorker

		// start the worker on its own goroutine
		go func() { _ = worker.Start() }()
		<-worker.NotifyStarted()
		q.AvailableWorkers <- worker
		q.Workers[x] = worker
	}
	q.Dispatch()
	return nil
}

// Dispatch processes work items in a loop.
func (q *Queue) Dispatch() {
	q.Latch.Started()
	defer q.Latch.Stopped()

	var workItem interface{}
	var worker *Worker
	var stopping <-chan struct{}
	for {
		stopping = q.Latch.NotifyStopping()
		select {
		case <-stopping:
			return
		case <-q.Background().Done():
			return
		default:
		}

		select {
		case <-stopping:
			return
		case <-q.Background().Done():
			return
		case workItem = <-q.Work:
			select {
			case <-stopping:
				q.Work <- workItem
				return
			case <-q.Background().Done():
				q.Work <- workItem
				return
			case worker = <-q.AvailableWorkers:
				worker.Enqueue(workItem)
			}
		}
	}
}

// Stop stops the queue and processes any remaining items.
func (q *Queue) Stop() error {
	if !q.Latch.CanStop() {
		return ex.New(ErrCannotStop)
	}
	q.Latch.WaitStopped() // wait for the dispatch loop to exit
	defer q.Latch.Reset() // reset the latch in case we have to start again

	timeoutContext, cancel := context.WithTimeout(q.Background(), q.ShutdownGracePeriod)
	defer cancel()

	if remainingWork := len(q.Work); remainingWork > 0 {
		for x := 0; x < remainingWork; x++ {
			// check the timeout first
			select {
			case <-timeoutContext.Done():
				return nil
			default:
			}

			select {
			case <-timeoutContext.Done():
				return nil
			case workItem := <-q.Work:
				select {
				case <-timeoutContext.Done():
					return nil
				case worker := <-q.AvailableWorkers:
					worker.Work <- workItem
				}
			}
		}
	}

	workersStopped := make(chan struct{})
	go func() {
		defer close(workersStopped)
		wg := sync.WaitGroup{}
		wg.Add(len(q.Workers))
		for _, worker := range q.Workers {
			go func(w *Worker) {
				defer wg.Done()
				w.StopContext(timeoutContext)
			}(worker)
		}
		wg.Wait()
	}()

	select {
	case <-timeoutContext.Done():
		return nil
	case <-workersStopped:
		return nil
	}
}

// Close stops the queue.
// Any work left in the queue will be discarded.
func (q *Queue) Close() error {
	q.Latch.WaitStopped()
	q.Latch.Reset()
	return nil
}

// ReturnWorker returns a given worker to the worker queue.
func (q *Queue) ReturnWorker(ctx context.Context, worker *Worker) error {
	q.AvailableWorkers <- worker
	return nil
}

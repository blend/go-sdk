/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package async

import (
	"context"

	"github.com/blend/go-sdk/ex"
)

// NewWorker creates a new worker.
func NewWorker(action WorkAction) *Worker {
	return &Worker{
		Latch:   NewLatch(),
		Context: context.Background(),
		Action:  action,
		Work:    make(chan interface{}),
	}
}

// Worker is a worker that is pushed work over a channel.
// It is used by other work distribution types (i.e. queue and batch)
// but can also be used independently.
type Worker struct {
	*Latch

	Context   context.Context
	Action    WorkAction
	Finalizer WorkerFinalizer

	SkipRecoverPanics bool
	Errors            chan error
	Work              chan interface{}
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
	return w.Latch.NotifyStopped()
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
	defer w.Latch.Stopped()

	var workItem interface{}
	var stopping <-chan struct{}
	for {
		stopping = w.Latch.NotifyStopping()
		select {
		case <-stopping:
			return
		case <-w.Background().Done():
			return
		default:
		}

		// block on work or stopping
		select {
		case workItem = <-w.Work:
			w.Execute(w.Background(), workItem)
		case <-stopping:
			return
		case <-w.Background().Done():
			return
		}
	}
}

// Execute invokes the action and recovers panics.
func (w *Worker) Execute(ctx context.Context, workItem interface{}) {
	defer func() {
		if !w.SkipRecoverPanics {
			if r := recover(); r != nil {
				w.HandleError(ex.New(r))
			}
		}
		if w.Finalizer != nil {
			w.HandleError(w.Finalizer(ctx, w))
		}
	}()
	if w.Action != nil {
		w.HandleError(w.Action(ctx, workItem))
	}
}

// Stop stops the worker.
//
// If there is an item left in the work channel
// it will be processed synchronously.
func (w *Worker) Stop() error {
	if !w.Latch.CanStop() {
		return ex.New(ErrCannotStop)
	}
	w.Latch.WaitStopped()
	w.Latch.Reset()
	return nil
}

// StopContext stops the worker in a given cancellation context.
func (w *Worker) StopContext(ctx context.Context) {
	stopped := make(chan struct{})
	go func() {
		defer func() {
			w.Latch.Reset()
			close(stopped)
		}()

		w.Latch.WaitStopped()
		if workLeft := len(w.Work); workLeft > 0 {
			for x := 0; x < workLeft; x++ {
				w.Execute(ctx, <-w.Work)
			}
		}
	}()
	select {
	case <-stopped:
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

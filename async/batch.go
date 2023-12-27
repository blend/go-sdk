/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package async

import (
	"context"
	"runtime"
)

// NewBatch creates a new batch processor.
// Batch processes are a known quantity of work that needs to be processed in parallel.
func NewBatch(work chan interface{}, action WorkAction, options ...BatchOption) *Batch {
	b := Batch{
		Action:      action,
		Work:        work,
		Parallelism: runtime.NumCPU(),
	}
	for _, option := range options {
		option(&b)
	}
	return &b
}

// BatchOption is an option for the batch worker.
type BatchOption func(*Batch)

// OptBatchErrors sets the batch worker error return channel.
func OptBatchErrors(errors chan error) BatchOption {
	return func(i *Batch) {
		i.Errors = errors
	}
}

// OptBatchSkipRecoverPanics sets the batch worker to throw (or to recover) panics.
func OptBatchSkipRecoverPanics(skipRecoverPanics bool) BatchOption {
	return func(i *Batch) {
		i.SkipRecoverPanics = skipRecoverPanics
	}
}

// OptBatchParallelism sets the batch worker parallelism, or the number of workers to create.
func OptBatchParallelism(parallelism int) BatchOption {
	return func(i *Batch) {
		i.Parallelism = parallelism
	}
}

// Batch is a batch of work executed by a fixed count of workers.
type Batch struct {
	Action            WorkAction
	SkipRecoverPanics bool
	Parallelism       int
	Work              chan interface{}
	Errors            chan error
}

// Process executes the action for all the work items.
func (b *Batch) Process(ctx context.Context) {
	if len(b.Work) == 0 {
		return
	}

	effectiveParallelism := b.Parallelism
	if effectiveParallelism == 0 {
		effectiveParallelism = runtime.NumCPU()
	}
	if effectiveParallelism > len(b.Work) {
		effectiveParallelism = len(b.Work)
	}

	allWorkers := make([]*Worker, effectiveParallelism)
	availableWorkers := make(chan *Worker, effectiveParallelism)

	// return worker is a local finalizer
	// that grabs a reference to the workers set.
	returnWorker := func(ctx context.Context, worker *Worker) error {
		availableWorkers <- worker
		return nil
	}

	// create and start workers.
	for x := 0; x < effectiveParallelism; x++ {
		worker := NewWorker(b.Action)
		worker.Context = ctx
		worker.Errors = b.Errors
		worker.Finalizer = returnWorker
		worker.SkipRecoverPanics = b.SkipRecoverPanics

		workerStarted := worker.NotifyStarted()
		go func() { _ = worker.Start() }()
		<-workerStarted

		allWorkers[x] = worker
		availableWorkers <- worker
	}
	defer func() {
		for x := 0; x < len(allWorkers); x++ {
			_ = allWorkers[x].Stop()
		}
	}()

	numWorkItems := len(b.Work)
	var worker *Worker
	var workItem interface{}
	for x := 0; x < numWorkItems; x++ {
		select {
		case <-ctx.Done():
			return
		default:
		}

		select {
		case workItem = <-b.Work:
			select {
			case worker = <-availableWorkers:
				select {
				case worker.Work <- workItem:
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

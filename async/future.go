/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package async

import (
	"context"
	"sync"
)

// Await returns a new future.
func Await(ctx context.Context, action ContextAction) *Future {
	f := &Future{
		action:   action,
		finished: make(chan error),
	}
	ctx, f.cancel = context.WithCancel(ctx)
	go f.do(ctx)
	return f
}

// Future is a deferred action.
type Future struct {
	mu       sync.Mutex
	action   ContextAction
	cancel   func()
	finished chan error
}

// Cancel quits the future.
func (f *Future) Cancel() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.cancel == nil {
		return ErrCannotCancel
	}
	f.cancel()
	err, ok := <-f.finished
	if ok {
		close(f.finished)
	}
	f.cancel = nil
	return err
}

// Complete blocks on the future completing.
func (f *Future) Complete() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	err, ok := <-f.finished
	if ok {
		close(f.finished)
		f.cancel = nil
	}
	return err
}

// Finished returns a channel that signals it is finished.
func (f *Future) Finished() <-chan error {
	return f.finished
}

func (f *Future) do(ctx context.Context) {
	f.finished <- f.action(ctx)
}

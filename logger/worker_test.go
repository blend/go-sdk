/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package logger

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestWorker(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(1)
	var didFire bool
	w := NewWorker(func(_ context.Context, e Event) {
		defer wg.Done()
		didFire = true

		typed, isTyped := e.(MessageEvent)
		assert.True(isTyped)
		assert.Equal("test", typed.Text)
	})

	go func() { _ = w.Start() }()
	<-w.NotifyStarted()
	defer func() { _ = w.Stop() }()

	w.Work <- EventWithContext{context.Background(), NewMessageEvent(Info, "test")}
	wg.Wait()

	assert.True(didFire)
}

func TestWorkerStop(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(4)
	var didFire bool
	w := NewWorker(func(ctx context.Context, e Event) {
		defer wg.Done()
		didFire = true
	})

	go func() { _ = w.Start() }()
	<-w.NotifyStarted()

	w.Work <- EventWithContext{Event: NewMessageEvent(Info, "test1")}
	w.Work <- EventWithContext{Event: NewMessageEvent(Info, "test2")}
	w.Work <- EventWithContext{Event: NewMessageEvent(Info, "test3")}
	w.Work <- EventWithContext{Event: NewMessageEvent(Info, "test4")}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		_ = w.StopContext(ctx)
	}()
	wg.Wait()

	assert.True(didFire)
}

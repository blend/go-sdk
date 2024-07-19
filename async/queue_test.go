/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package async

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_Queue_Start_Close(t *testing.T) {
	its := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(8)
	q := NewQueue(func(_ context.Context, obj interface{}) error {
		defer wg.Done()
		return nil
	})

	go func() { _ = q.Start() }()
	<-q.Latch.NotifyStarted()

	its.True(q.Latch.IsStarted())

	for x := 0; x < 8; x++ {
		q.Enqueue(fmt.Sprint(x))
	}

	wg.Wait()
	q.Close()
	its.False(q.Latch.IsStarted())
}

func Test_Queue_Start_Stop(t *testing.T) {
	its := assert.New(t)

	workCount := 128

	wg := sync.WaitGroup{}
	wg.Add(workCount)
	q := NewQueue(func(_ context.Context, obj interface{}) error {
		defer wg.Done()
		return nil
	}, OptQueueMaxWork(workCount+1))

	go func() { _ = q.Start() }()
	<-q.Latch.NotifyStarted()
	its.True(q.Latch.IsStarted())

	for x := 0; x < workCount; x++ {
		q.Enqueue(fmt.Sprint(x))
	}

	its.Nil(q.Stop())
	its.False(q.Latch.IsStarted())
	wg.Wait()
}

func Test_Queue_Start_Stop_Start(t *testing.T) {
	its := assert.New(t)

	workCount := 10

	wg := sync.WaitGroup{}
	wg.Add(workCount)
	q := NewQueue(func(_ context.Context, obj interface{}) error {
		defer wg.Done()
		return nil
	}, OptQueueMaxWork(workCount))

	go func() { _ = q.Start() }()
	<-q.Latch.NotifyStarted()

	its.True(q.Latch.IsStarted())

	for x := 0; x < workCount; x++ {
		q.Enqueue(fmt.Sprint(x))
	}
	its.Nil(q.Stop())
	its.False(q.Latch.IsStarted())
	wg.Wait()

	wg.Add(workCount)

	go func() { _ = q.Start() }()
	<-q.Latch.NotifyStarted()

	its.True(q.Latch.IsStarted())

	for x := 0; x < workCount; x++ {
		q.Enqueue(fmt.Sprint(x))
	}
	its.Nil(q.Stop())
	its.False(q.Latch.IsStarted())
	wg.Wait()
}

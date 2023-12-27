/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package async

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/graceful"
)

// Assert a latch is graceful
var (
	_ graceful.Graceful = (*Interval)(nil)
)

func Test_Interval(t *testing.T) {
	assert := assert.New(t)

	var didWork bool
	unbuffered := make(chan bool)
	w := NewInterval(func(_ context.Context) error {
		didWork = true
		<-unbuffered
		return nil
	}, time.Millisecond)

	assert.Equal(time.Millisecond, w.Interval)

	go func() { _ = w.Start() }()
	<-w.NotifyStarted()

	assert.True(w.IsStarted())
	unbuffered <- true
	close(unbuffered)
	assert.Nil(w.Stop())
	assert.True(w.IsStopped())
	assert.True(didWork)
}

func Test_Interval_StopOnError(t *testing.T) {
	its := assert.New(t)

	var didWork bool
	unbuffered := make(chan bool)
	w := NewInterval(func(_ context.Context) error {
		didWork = true
		<-unbuffered
		return fmt.Errorf("this is just a test")
	}, time.Millisecond, OptIntervalStopOnError(true))

	its.Equal(time.Millisecond, w.Interval)

	startErrors := make(chan error)
	go func() {
		startErrors <- w.Start()
	}()
	<-w.NotifyStarted()

	its.True(w.IsStarted())
	unbuffered <- true
	close(unbuffered)
	its.True(didWork)
	err := <-startErrors
	its.Equal(fmt.Errorf("this is just a test"), err)
	its.True(w.IsStopped())
}

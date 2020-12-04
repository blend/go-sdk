package async

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestBatch(t *testing.T) {
	assert := assert.New(t)

	workItems := 32

	items := make(chan interface{}, workItems)
	for x := 0; x < workItems; x++ {
		items <- "hello" + strconv.Itoa(x)
	}

	var processed int32
	action := func(_ context.Context, v interface{}) error {
		atomic.AddInt32(&processed, 1)
		return fmt.Errorf("this is only a test")
	}

	errors := make(chan error, workItems)
	NewBatch(
		items,
		action,
		OptBatchErrors(errors),
		OptBatchParallelism(4),
	).Process(context.Background())

	assert.Equal(workItems, processed)
	assert.Equal(workItems, len(errors))
}

func TestBatchPanic(t *testing.T) {
	assert := assert.New(t)

	workItems := 32

	items := make(chan interface{}, workItems)
	for x := 0; x < workItems; x++ {
		items <- "hello" + strconv.Itoa(x)
	}

	var processed int32
	action := func(_ context.Context, v interface{}) error {
		if result := atomic.AddInt32(&processed, 1); result == 1 {
			panic("this is only a test")
		}
		return nil
	}

	errors := make(chan error, workItems)
	NewBatch(items, action, OptBatchErrors(errors)).Process(context.Background())

	assert.Equal(workItems, processed)
	assert.Equal(1, len(errors))
}

func TestBatchCancel(t *testing.T) {
	assert := assert.New(t)

	workItems := 32

	items := make(chan interface{}, workItems)
	for x := 0; x < workItems; x++ {
		items <- "hello" + strconv.Itoa(x)
	}

	var processed int32
	var didCancel bool
	countReached := make(chan struct{})
	action := func(ctx context.Context, _ interface{}) error {
		if result := atomic.AddInt32(&processed, 1); int(result) > workItems>>1 {
			if !didCancel {
				close(countReached)
				select {
				case <-ctx.Done():
					didCancel = true
				}
			}
		}
		return nil
	}

	errors := make(chan error, workItems)
	withCancel, cancel := context.WithCancel(context.Background())

	go func() {
		<-countReached
		cancel()
	}()

	NewBatch(items, action, OptBatchErrors(errors)).Process(withCancel)
	// assert.True(int32(workItems) > processed)
	assert.True(didCancel)
}

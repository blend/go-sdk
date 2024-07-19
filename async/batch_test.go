/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package async

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_Batch(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

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

	its.Equal(workItems, processed)
	its.Equal(workItems, len(errors))
}

func Test_Batch_empty(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	items := make(chan interface{}, 32)

	var processed int32
	action := func(_ context.Context, v interface{}) error {
		atomic.AddInt32(&processed, 1)
		return fmt.Errorf("this is only a test")
	}

	errors := make(chan error, 32)
	NewBatch(
		items,
		action,
		OptBatchErrors(errors),
		OptBatchParallelism(4),
	).Process(context.Background())

	its.Equal(0, processed)
	its.Equal(0, len(errors))
}

func Test_Batch_panic(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

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

	its.Equal(workItems, processed)
	its.Equal(1, len(errors))
}

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

	items := make(chan interface{}, 32)
	for x := 0; x < 32; x++ {
		items <- "hello" + strconv.Itoa(x)
	}

	var processed int32
	action := func(_ context.Context, v interface{}) error {
		atomic.AddInt32(&processed, 1)
		return fmt.Errorf("this is only a test")
	}

	errors := make(chan error, 32)
	NewBatch(action, items, OptBatchErrors(errors)).Process(context.Background())

	assert.Equal(32, processed)
	assert.Equal(32, len(errors))
}

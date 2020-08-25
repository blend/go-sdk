package async

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/graceful"
)

// Assert a latch is graceful
var (
	_ graceful.Graceful = (*Interval)(nil)
)

func TestIntervalWorker(t *testing.T) {
	assert := assert.New(t)

	var didWork bool
	var didWorkLock sync.Mutex
	done := sync.WaitGroup{}
	done.Add(1)
	wait := make(chan struct{})
	w := NewInterval(func(_ context.Context) error {
		didWorkLock.Lock()
		defer didWorkLock.Unlock()
		if didWork {
			return nil
		}

		defer done.Done()
		<-wait
		didWork = true
		return nil
	}, time.Millisecond)

	assert.Equal(time.Millisecond, w.Interval)

	go w.Start()
	<-w.NotifyStarted()

	assert.True(w.IsStarted())
	close(wait)
	done.Wait()
	w.Stop()
	assert.True(w.IsStopped())
	assert.True(didWork)
}

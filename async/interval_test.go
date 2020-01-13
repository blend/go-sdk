package async

import (
	"context"
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
	done := make(chan struct{})
	w := NewInterval(func(_ context.Context) error {
		defer func() {
			close(done)
		}()
		didWork = true
		return nil
	}, time.Millisecond)

	assert.Equal(time.Millisecond, w.Interval)

	go w.Start()
	<-w.NotifyStarted()

	assert.True(w.IsStarted())
	<-done
	w.Stop()
	assert.True(w.IsStopped())
	assert.True(didWork)
}

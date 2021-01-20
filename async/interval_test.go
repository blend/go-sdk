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

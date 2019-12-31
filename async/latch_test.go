package async

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestLatch(t *testing.T) {
	assert := assert.New(t)

	l := NewLatch()

	var didStart bool
	var didAbort bool
	var didGetWork bool

	work := make(chan bool)
	workComplete := make(chan bool)

	l.Starting()
	assert.True(l.IsStarting())
	assert.False(l.IsStarted())
	assert.False(l.IsStopping())
	assert.False(l.IsStopped())

	go func() {
		l.Started()
		didStart = true
		for {
			select {
			case <-work:
				didGetWork = true
				workComplete <- true
			case <-l.NotifyStopping():
				didAbort = true
				l.Stopped()
				return
			}
		}
	}()
	<-l.NotifyStarted()

	work <- true
	assert.True(l.IsStarted())

	// wait for work to happen.
	<-workComplete

	// signal stop
	l.Stopping()
	<-l.NotifyStopped()

	assert.True(didStart)
	assert.True(didAbort)
	assert.True(didGetWork)
	assert.False(l.IsStopping())
	assert.False(l.IsStarted())
	assert.True(l.IsStopped())

	didAbort = false
	assert.False(didAbort)
	assert.False(l.IsStopping())
	assert.False(l.IsStarted())
	assert.True(l.IsStopped())
}

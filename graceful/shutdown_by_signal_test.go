package graceful

import (
	"fmt"
	"os"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func newHosted() *hosted {
	return &hosted{
		started: make(chan struct{}),
		stopped: make(chan struct{}),
	}
}

type hosted struct {
	state   int32
	started chan struct{}
	stopped chan struct{}
}

func (h *hosted) Start() error {
	h.state = 1
	h.stopped = make(chan struct{})
	close(h.started)
	<-h.stopped
	return nil
}

func (h *hosted) Stop() error {
	if h.state != 1 {
		return fmt.Errorf("cannot stop")
	}
	h.state = 0
	h.started = make(chan struct{})
	close(h.stopped)
	return nil
}

func (h *hosted) NotifyStarted() <-chan struct{} {
	return h.started
}

func (h *hosted) NotifyStopped() <-chan struct{} {
	return h.stopped
}

func TestShutdownBySignal(t *testing.T) {
	assert := assert.New(t)

	hosted := newHosted()

	terminateSignal := make(chan os.Signal)
	var err error
	done := make(chan struct{})
	go func() {
		err = ShutdownBySignal([]Graceful{hosted}, OptShutdownSignal(terminateSignal))
		close(done)
	}()
	<-hosted.NotifyStarted()

	close(terminateSignal)
	<-done
	assert.Nil(err)
}

func TestShutdownBySignalMany(t *testing.T) {
	assert := assert.New(t)

	workers := []Graceful{
		newHosted(),
		newHosted(),
		newHosted(),
		newHosted(),
		newHosted(),
	}

	terminateSignal := make(chan os.Signal)
	var err error
	done := make(chan struct{})

	go func() {
		err = ShutdownBySignal(workers, OptShutdownSignal(terminateSignal))
		close(done)
	}()

	// wait for the workers to start
	for _, h := range workers {
		<-h.(*hosted).started
	}
	for _, h := range workers {
		assert.Equal(1, h.(*hosted).state)
	}

	close(terminateSignal)
	<-done
	assert.Nil(err)
	for _, h := range workers {
		assert.Equal(0, h.(*hosted).state)
	}
}

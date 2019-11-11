package logger

import (
	"bytes"
	"sync"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestWorker(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(1)
	var didFire bool
	w := NewWorker(nil, func(e Event) {
		defer wg.Done()
		didFire = true

		typed, isTyped := e.(*MessageEvent)
		assert.True(isTyped)
		assert.Equal("test", typed.Message())
	}, DefaultWriteQueueDepth)

	w.Start()
	defer w.Close()

	w.Work <- Messagef(Info, "test")
	wg.Wait()

	assert.True(didFire)
}

func TestWorkerPanics(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer(nil)
	wr := NewTextWriter(buffer)

	log := New().WithFlags(AllFlags()).WithWriter(wr)
	defer log.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)
	var didFire bool
	w := NewWorker(log, func(e Event) {
		defer wg.Done()
		didFire = true
		panic("only a test")
	}, DefaultWriteQueueDepth)
	w.Start()

	w.Work <- Messagef(Info, "test")
	wg.Wait()

	assert.True(didFire)
	w.Close()
	assert.NotEmpty(buffer.String())
}

func TestWorkerDrain(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(4)
	var didFire bool
	w := NewWorker(nil, func(e Event) {
		defer wg.Done()
		didFire = true
	}, DefaultWriteQueueDepth)

	w.Work <- Messagef(Info, "test1")
	w.Work <- Messagef(Info, "test2")
	w.Work <- Messagef(Info, "test3")
	w.Work <- Messagef(Info, "test4")

	go func() {
		w.Drain()
	}()
	wg.Wait()

	assert.True(didFire)
}

package bufferutil

import (
	"context"
	"sync"

	"github.com/blend/go-sdk/async"
)

// BufferHandlers is a synchronized map of listeners for new lines to a line buffer.
type BufferHandlers struct {
	sync.RWMutex
	Handlers map[string]*async.Worker
}

// Add adds a handler.
func (bl *BufferHandlers) Add(uid string, handler BufferChunkHandler) {
	bl.Lock()
	defer bl.Unlock()

	if bl.Handlers == nil {
		bl.Handlers = make(map[string]*async.Worker)
	}

	w := async.NewWorker(func(_ context.Context, wi interface{}) error {
		handler(wi.(BufferChunk))
		return nil
	})
	w.Work = make(chan interface{}, 128)

	go w.Start()

	bl.Handlers[uid] = w
}

// Remove removes a listener.
func (bl *BufferHandlers) Remove(uid string) {
	bl.Lock()
	defer bl.Unlock()

	if bl.Handlers == nil {
		bl.Handlers = make(map[string]*async.Worker)
	}
	if w, ok := bl.Handlers[uid]; ok {
		stopped := w.NotifyStopped()
		w.Stop()
		<-stopped
	}
	delete(bl.Handlers, uid)
}

// Handle calls the handlers.
func (bl *BufferHandlers) Handle(chunk BufferChunk) {
	bl.RLock()
	defer bl.RUnlock()

	for _, handler := range bl.Handlers {
		handler.Work <- chunk
	}
}

// Close shuts down the handlers.
func (bl *BufferHandlers) Close() {
	bl.Lock()
	defer bl.Unlock()

	for _, queue := range bl.Handlers {
		stopped := queue.NotifyStopped()
		queue.Stop()
		<-stopped
	}
}

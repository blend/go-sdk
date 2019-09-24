package cron

import (
	"context"
	"sync"

	"github.com/blend/go-sdk/async"
)

// OutputListeners is a synchronized map of listeners for new lines to a line buffer.
type OutputListeners struct {
	sync.RWMutex
	Listeners map[string]*async.Queue
}

// Add adds a listener.
func (ol *OutputListeners) Add(uid string, listener OutputListener) {
	ol.Lock()
	if ol.Listeners == nil {
		ol.Listeners = make(map[string]*async.Queue)
	}
	w := async.NewQueue(func(_ context.Context, wi interface{}) error {
		listener(wi.(OutputChunk))
		return nil
	}, async.OptQueueMaxWork(128))
	go w.Start()

	ol.Listeners[uid] = w
	ol.Unlock()
}

// Remove removes a listener.
func (ol *OutputListeners) Remove(uid string) {
	ol.Lock()
	if ol.Listeners == nil {
		ol.Listeners = make(map[string]*async.Queue)
	}
	if w, ok := ol.Listeners[uid]; ok {
		w.Stop()
	}
	delete(ol.Listeners, uid)
	ol.Unlock()
}

// Trigger calls the handlers.
func (ol *OutputListeners) Trigger(line OutputChunk) {
	ol.RLock()
	defer ol.RUnlock()

	for index := range ol.Listeners {
		ol.Listeners[index].Work <- line
	}
}

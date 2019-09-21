package cron

import "sync"

// OutputListeners is a synchronized map of listeners for new lines to a line buffer.
type OutputListeners struct {
	sync.RWMutex
	Listeners map[string]OutputListener
}

// Add adds a listener.
func (ol *OutputListeners) Add(uid string, listener OutputListener) {
	ol.Lock()
	if ol.Listeners == nil {
		ol.Listeners = make(map[string]OutputListener)
	}
	ol.Listeners[uid] = listener
	ol.Unlock()
}

// Remove removes a listener.
func (ol *OutputListeners) Remove(uid string) {
	ol.Lock()
	if ol.Listeners == nil {
		ol.Listeners = make(map[string]OutputListener)
	}
	delete(ol.Listeners, uid)
	ol.Unlock()
}

// Trigger calls the handlers.
func (ol *OutputListeners) Trigger(line OutputChunk) {
	ol.RLock()
	defer ol.RUnlock()

	for _, listener := range ol.Listeners {
		listener(line)
	}
}

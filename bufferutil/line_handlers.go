package bufferutil

import "sync"

// LineHandlers is a synchronized map of listeners for new lines to a line buffer.
type LineHandlers struct {
	sync.RWMutex
	Handlers map[string]LineHandler
}

// Add adds a listener.
func (lh *LineHandlers) Add(uid string, handler LineHandler) {
	lh.Lock()
	if lh.Handlers == nil {
		lh.Handlers = make(map[string]LineHandler)
	}
	lh.Handlers[uid] = handler
	lh.Unlock()
}

// Remove removes a listener.
func (lh *LineHandlers) Remove(uid string) {
	lh.Lock()
	if lh.Handlers == nil {
		lh.Handlers = make(map[string]LineHandler)
	}
	delete(lh.Handlers, uid)
	lh.Unlock()
}

// Handle calls the handlers.
func (lh *LineHandlers) Handle(line Line) {
	lh.RLock()
	defer lh.RUnlock()
	for _, handler := range lh.Handlers {
		handler(line)
	}
}

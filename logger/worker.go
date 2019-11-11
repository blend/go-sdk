package logger

import (
	"sync"
)

// NewWorker returns a new worker.
func NewWorker(parent *Logger, listener Listener, queueDepth int) *Worker {
	return &Worker{
		Parent:   parent,
		Listener: listener,
		Work:     make(chan Event, queueDepth),
	}
}

// Worker is an agent that processes a listener.
type Worker struct {
	sync.Mutex
	Parent   *Logger
	Listener Listener
	Abort    chan struct{}
	Aborted  chan struct{}
	Work     chan Event
}

// Start starts the worker.
func (w *Worker) Start() {
	w.Lock()
	w.startUnsafe()
	w.Unlock()
}

func (w *Worker) startUnsafe() {
	w.Abort = make(chan struct{})
	w.Aborted = make(chan struct{})
	go w.ProcessLoop()
}

// ProcessLoop is the for/select loop.
func (w *Worker) ProcessLoop() {
	var e Event
	for {
		select {
		case e = <-w.Work:
			w.Process(e)
		case <-w.Abort:
			close(w.Aborted)
			return
		}
	}
}

// Process calls the listener for an event.
func (w *Worker) Process(e Event) {
	if w.Parent != nil && w.Parent.RecoversPanics() {
		defer func() {
			if r := recover(); r != nil {
				if w.Parent != nil {
					w.Parent.Write(Errorf(Fatal, "%+v", r))
				}
			}
		}()
	}
	w.Listener(e)
}

// Drain stops the worker and synchronously processes any remaining work.
// It then restarts the worker.
func (w *Worker) Drain() {
	w.Lock()
	defer w.Unlock()

	if w.Abort != nil {
		close(w.Abort)
		<-w.Aborted
	}

	workLeft := len(w.Work)
	for index := 0; index < workLeft; index++ {
		w.Process(<-w.Work)
	}

	w.startUnsafe()
}

// Close closes the worker.
func (w *Worker) Close() error {
	w.Lock()
	defer w.Unlock()

	close(w.Abort)
	<-w.Aborted

	for len(w.Work) > 0 {
		w.Process(<-w.Work)
	}
	close(w.Work)

	return nil
}

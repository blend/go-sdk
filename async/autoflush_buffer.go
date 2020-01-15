package async

import (
	"context"
	"sync"
	"time"

	"github.com/blend/go-sdk/collections"
	"github.com/blend/go-sdk/ex"
)

// NewAutoflushBuffer creates a new autoflush buffer.
func NewAutoflushBuffer(handler AutoflushAction, options ...AutoflushBufferOption) *AutoflushBuffer {
	afb := AutoflushBuffer{
		Latch:       NewLatch(),
		Handler:     handler,
		MaxLen:      DefaultQueueMaxWork,
		Interval:    DefaultInterval,
		FlushOnStop: true,
	}
	for _, option := range options {
		option(&afb)
	}
	return &afb
}

// AutoflushBufferOption is an option for auto-flush buffers.
type AutoflushBufferOption func(*AutoflushBuffer)

// OptAutoflushBufferMaxLen sets the auto-flush buffer's maximum length.
func OptAutoflushBufferMaxLen(maxLen int) AutoflushBufferOption {
	return func(afb *AutoflushBuffer) {
		afb.MaxLen = maxLen
	}
}

// OptAutoflushBufferInterval sets the auto-flush buffer's interval.
func OptAutoflushBufferInterval(d time.Duration) AutoflushBufferOption {
	return func(afb *AutoflushBuffer) {
		afb.Interval = d
	}
}

// OptAutoflushBufferContext sets the auto-flush buffer's context.
func OptAutoflushBufferContext(ctx context.Context) AutoflushBufferOption {
	return func(afb *AutoflushBuffer) {
		afb.Context = ctx
	}
}

// OptAutoflushBufferErrors sets the auto-flush buffer's error return channel.
func OptAutoflushBufferErrors(errors chan error) AutoflushBufferOption {
	return func(afb *AutoflushBuffer) {
		afb.Errors = errors
	}
}

// OptAutoflushBufferFlushOnStop sets the auto-flush buffer's flush on stop option.
func OptAutoflushBufferFlushOnStop(flushOnStop bool) AutoflushBufferOption {
	return func(afb *AutoflushBuffer) {
		afb.FlushOnStop = flushOnStop
	}
}

// AutoflushBuffer is a backing store that operates either on a fixed length flush or a fixed interval flush.
// A handler should be provided but without one the buffer will just clear.
// Adds that would cause fixed length flushes do not block on the flush handler.
type AutoflushBuffer struct {
	sync.Mutex
	Latch       *Latch
	Context     context.Context
	MaxLen      int
	Interval    time.Duration
	Contents    *collections.RingBuffer
	FlushOnStop bool
	Handler     AutoflushAction
	Errors      chan error
}

// Background returns a background context.
func (ab *AutoflushBuffer) Background() context.Context {
	if ab.Context != nil {
		return ab.Context
	}
	return context.Background()
}

/*
Start starts the auto-flush buffer.
This call blocks. To call it asynchronously:

	go afb.Start()
	<-afb.NotifyStarted()
*/
func (ab *AutoflushBuffer) Start() error {
	if !ab.Latch.CanStart() {
		return ex.New(ErrCannotStart)
	}
	ab.Latch.Starting()
	ab.Contents = collections.NewRingBufferWithCapacity(ab.MaxLen)
	ab.Dispatch()
	return nil
}

// Dispatch is the main run loop.
func (ab *AutoflushBuffer) Dispatch() {
	ab.Latch.Started()
	ticker := time.Tick(ab.Interval)
	var stopping <-chan struct{}
	for {
		stopping = ab.Latch.NotifyStopping()
		select {
		case <-ticker:
			ab.FlushAsync(ab.Background())
		case <-stopping:
			if ab.FlushOnStop {
				ab.Flush(ab.Background())
			}
			ab.Latch.Stopped()
			return
		}
	}
}

// Stop stops the buffer flusher.
func (ab *AutoflushBuffer) Stop() error {
	if !ab.Latch.CanStop() {
		return ex.New(ErrCannotStop)
	}
	ab.Latch.Stopping()
	<-ab.Latch.NotifyStopped()
	return nil
}

// NotifyStarted implements graceful.Graceful.
func (ab *AutoflushBuffer) NotifyStarted() <-chan struct{} {
	return ab.Latch.NotifyStarted()
}

// NotifyStopped implements graceful.Graceful.
func (ab *AutoflushBuffer) NotifyStopped() <-chan struct{} {
	return ab.Latch.NotifyStopped()
}

// Add adds a new object to the buffer, blocking if it triggers a flush.
// If the buffer is full, it will call the flush handler on a separate goroutine.
func (ab *AutoflushBuffer) Add(obj interface{}) {
	ab.Lock()
	defer ab.Unlock()

	ab.Contents.Enqueue(obj)
	if ab.Contents.Len() >= ab.MaxLen {
		ab.flushUnsafeAsync(ab.Background(), ab.Contents.Drain())
	}
}

// AddMany adds many objects to the buffer at once.
func (ab *AutoflushBuffer) AddMany(objs ...interface{}) {
	ab.Lock()
	defer ab.Unlock()

	for _, obj := range objs {
		ab.Contents.Enqueue(obj)
		if ab.Contents.Len() >= ab.MaxLen {
			ab.flushUnsafeAsync(ab.Background(), ab.Contents.Drain())
		}
	}
}

// Flush clears the buffer, if a handler is provided it is passed the contents of the buffer.
// This call is synchronous, in that it will call the flush handler on the same goroutine.
func (ab *AutoflushBuffer) Flush(ctx context.Context) {
	ab.Lock()
	defer ab.Unlock()
	ab.flushUnsafe(ctx, ab.Contents.Drain())
}

// FlushAsync clears the buffer, if a handler is provided it is passed the contents of the buffer.
// This call is asynchronous, in that it will call the flush handler on its own goroutine.
func (ab *AutoflushBuffer) FlushAsync(ctx context.Context) {
	ab.Lock()
	defer ab.Unlock()
	ab.flushUnsafeAsync(ctx, ab.Contents.Drain())
}

// flushUnsafeAsync flushes the buffer without acquiring any locks.
func (ab *AutoflushBuffer) flushUnsafeAsync(ctx context.Context, contents []interface{}) {
	go ab.flushUnsafe(ctx, contents)
}

// flushUnsafeAsync flushes the buffer without acquiring any locks.
func (ab *AutoflushBuffer) flushUnsafe(ctx context.Context, contents []interface{}) {
	if ab.Handler != nil {
		if len(contents) > 0 {
			Recover(func() error {
				return ab.Handler(ctx, contents)
			}, ab.Errors)
		}
	}
}

package autoflush

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/collections"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stats"
	"github.com/blend/go-sdk/timeutil"
)

// New creates a new  buffer.
func New(handler Action, options ...Option) *Buffer {
	afb := Buffer{
		Latch:               async.NewLatch(),
		Handler:             handler,
		Parallelism:         runtime.NumCPU(),
		MaxFlushes:          DefaultMaxFlushes,
		MaxLen:              DefaultMaxLen,
		Interval:            DefaultFlushInterval,
		ShutdownGracePeriod: DefaultShutdownGracePeriod,
	}
	for _, option := range options {
		option(&afb)
	}
	afb.contents = collections.NewRingBufferWithCapacity(afb.MaxLen)
	return &afb
}

// Option is an option for auto-flush buffers.
type Option func(*Buffer)

// OptMaxFlushes sets the auto-flush buffer's maximum flush queue length.
func OptMaxFlushes(maxFlushes int) Option {
	return func(afb *Buffer) {
		afb.MaxFlushes = maxFlushes
	}
}

// OptParallelism sets the auto-flush buffer's flush worker count.
func OptParallelism(parallelism int) Option {
	return func(afb *Buffer) {
		afb.Parallelism = parallelism
	}
}

// OptMaxLen sets the auto-flush buffer's maximum length.
func OptMaxLen(maxLen int) Option {
	return func(afb *Buffer) {
		afb.MaxLen = maxLen
	}
}

// OptInterval sets the auto-flush buffer's interval.
func OptInterval(d time.Duration) Option {
	return func(afb *Buffer) {
		afb.Interval = d
	}
}

// OptContext sets the auto-flush buffer's context.
func OptContext(ctx context.Context) Option {
	return func(afb *Buffer) {
		afb.Context = ctx
	}
}

// OptErrors sets the auto-flush buffer's error return channel.
func OptErrors(errors chan error) Option {
	return func(afb *Buffer) {
		afb.Errors = errors
	}
}

// OptShutdownGracePeriod sets the auto-flush buffer's shutdown grace period.
func OptShutdownGracePeriod(shutdownGracePeriod time.Duration) Option {
	return func(afb *Buffer) {
		afb.ShutdownGracePeriod = shutdownGracePeriod
	}
}

// OptLog sets the Buffer logger.
func OptLog(log logger.Log) Option {
	return func(afb *Buffer) {
		afb.Log = log
	}
}

// OptStats sets the Buffer stats collector.
func OptStats(stats stats.Collector) Option {
	return func(afb *Buffer) {
		afb.Stats = stats
	}
}

// OptTracer sets the Buffer logger.
func OptTracer(tracer Tracer) Option {
	return func(afb *Buffer) {
		afb.Tracer = tracer
	}
}

// Action is an action called by an  buffer.
type Action func(context.Context, []interface{}) error

// Buffer is a backing store that operates either on a fixed length flush or a fixed interval flush.
// A handler should be provided but without one the buffer will just clear.
// Adds that would cause fixed length flushes do not block on the flush handler.
type Buffer struct {
	Latch   *async.Latch
	Context context.Context

	Log    logger.Log
	Stats  stats.Collector
	Tracer Tracer

	MaxLen              int
	Interval            time.Duration
	Parallelism         int
	MaxFlushes          int
	ShutdownGracePeriod time.Duration

	contentsMu sync.Mutex
	contents   *collections.RingBuffer

	Handler Action
	Errors  chan error

	intervalWorker    *async.Interval
	flushes           chan Flush
	flushWorkersReady chan *async.Worker
	flushWorkers      []*async.Worker
}

// Background returns a background context.
func (ab *Buffer) Background() context.Context {
	if ab.Context != nil {
		return ab.Context
	}
	return context.Background()
}

//Start starts the auto-flush buffer.
/*
This call blocks. To call it asynchronously:

	go afb.Start()
	<-afb.NotifyStarted()
*/
func (ab *Buffer) Start() error {
	if !ab.Latch.CanStart() {
		return ex.New(async.ErrCannotStart)
	}
	ab.Latch.Starting()

	ab.flushes = make(chan Flush, ab.MaxFlushes)
	ab.flushWorkers = make([]*async.Worker, ab.Parallelism)
	ab.flushWorkersReady = make(chan *async.Worker, ab.Parallelism)
	ab.intervalWorker = async.NewInterval(ab.FlushAsync, ab.Interval, async.OptIntervalErrors(ab.Errors))

	for x := 0; x < ab.Parallelism; x++ {
		worker := async.NewWorker(ab.workerAction)
		worker.Context = ab.Context
		worker.Errors = ab.Errors
		worker.Finalizer = ab.returnFlushWorker
		go func() { _ = worker.Start() }()
		<-worker.NotifyStarted()
		ab.flushWorkers[x] = worker
		ab.flushWorkersReady <- worker
	}
	go func() { _ = ab.intervalWorker.Start() }()
	ab.Dispatch()
	return nil
}

// Dispatch is the main run loop.
func (ab *Buffer) Dispatch() {
	ab.Latch.Started()

	var stopping <-chan struct{}
	var flushWorker *async.Worker
	var flush Flush
	for {
		stopping = ab.Latch.NotifyStopping()
		select {
		case <-stopping:
			ab.Latch.Stopped()
			return
		default:
		}
		select {
		case flush = <-ab.flushes:
			select {
			case flushWorker = <-ab.flushWorkersReady:
				flushWorker.Work <- flush
			case <-stopping:
				ab.flushes <- flush
				ab.Latch.Stopped()
				return
			}
		case <-stopping:
			ab.Latch.Stopped()
			return
		}
	}
}

// Stop stops the buffer flusher.
//
// Any in flight flushes will be given ShutdownGracePeriod amount of time.
//
// Stop is _very_ complicated.
func (ab *Buffer) Stop() error {
	if !ab.Latch.CanStop() {
		return ex.New(async.ErrCannotStop)
	}
	// stop the interval worker
	ab.intervalWorker.WaitStopped()

	// stop the running dispatch loop
	stopped := ab.Latch.NotifyStopped()
	ab.Latch.Stopping()
	<-stopped

	timeoutContext, cancel := context.WithTimeout(ab.Background(), ab.ShutdownGracePeriod)
	defer func() {
		cancel()
	}()

	ab.contentsMu.Lock()
	defer ab.contentsMu.Unlock()
	if ab.contents.Len() > 0 {
		ab.flushes <- Flush{
			Context:  timeoutContext,
			Contents: ab.contents.Drain(),
		}
	}

	if remainingFlushes := len(ab.flushes); remainingFlushes > 0 {
		logger.MaybeDebugf(ab.Log, "%d flushes remaining", remainingFlushes)
		var flushWorker *async.Worker
		var flush Flush
		for x := 0; x < remainingFlushes; x++ {
			select {
			case <-timeoutContext.Done():
				logger.MaybeDebugf(ab.Log, "stop timed out")
				return nil
			case flush = <-ab.flushes:
				select {
				case <-timeoutContext.Done():
					logger.MaybeDebugf(ab.Log, "stop timed out")
					return nil
				case flushWorker = <-ab.flushWorkersReady:
					flushWorker.Work <- flush
				}
			}
		}
	}

	workersStopped := make(chan struct{})
	go func() {
		defer close(workersStopped)
		wg := sync.WaitGroup{}
		wg.Add(len(ab.flushWorkers))
		for index, worker := range ab.flushWorkers {
			go func(i int, w *async.Worker) {
				defer wg.Done()
				logger.MaybeDebugf(ab.Log, "draining worker %d", i)
				w.Drain(timeoutContext)
			}(index, worker)
		}
		wg.Wait()
	}()

	select {
	case <-timeoutContext.Done():
		logger.MaybeDebugf(ab.Log, "stop timed out")
		return nil
	case <-workersStopped:
		return nil
	}
}

// NotifyStarted implements graceful.Graceful.
func (ab *Buffer) NotifyStarted() <-chan struct{} {
	return ab.Latch.NotifyStarted()
}

// NotifyStopped implements graceful.Graceful.
func (ab *Buffer) NotifyStopped() <-chan struct{} {
	return ab.Latch.NotifyStopped()
}

// Add adds a new object to the buffer, blocking if it triggers a flush.
// If the buffer is full, it will call the flush handler on a separate goroutine.
func (ab *Buffer) Add(ctx context.Context, obj interface{}) {
	if ab.Tracer != nil {
		finisher := ab.Tracer.StartAdd(ctx)
		defer finisher.Finish(nil)
	}
	var bufferLength int
	if ab.Stats != nil {
		ab.maybeStatCount(ctx, MetricAdd, 1)
		start := time.Now().UTC()
		defer func() {
			ab.maybeStatGauge(ctx, MetricBufferLength, float64(bufferLength))
			ab.maybeStatElapsed(ctx, MetricAddElapsed, start)
		}()
	}

	var flush []interface{}
	ab.contentsMu.Lock()
	bufferLength = ab.contents.Len()
	ab.contents.Enqueue(obj)
	if ab.contents.Len() >= ab.MaxLen {
		flush = ab.contents.Drain()
	}
	ab.contentsMu.Unlock()
	ab.unsafeFlushAsync(ctx, flush)
}

// AddMany adds many objects to the buffer at once.
func (ab *Buffer) AddMany(ctx context.Context, objs ...interface{}) {
	if ab.Tracer != nil {
		finisher := ab.Tracer.StartAddMany(ctx)
		defer finisher.Finish(nil)
	}
	var bufferLength int
	if ab.Stats != nil {
		ab.maybeStatCount(ctx, MetricAddMany, 1)
		ab.maybeStatCount(ctx, MetricAddManyItemCount, len(objs))
		start := time.Now().UTC()
		defer func() {
			ab.maybeStatGauge(ctx, MetricBufferLength, float64(bufferLength))
			ab.maybeStatElapsed(ctx, MetricAddManyElapsed, start)
		}()
	}

	var flushes [][]interface{}
	ab.contentsMu.Lock()
	bufferLength = ab.contents.Len()
	for _, obj := range objs {
		ab.contents.Enqueue(obj)
		if ab.contents.Len() >= ab.MaxLen {
			flushes = append(flushes, ab.contents.Drain())
		}
	}
	ab.contentsMu.Unlock()
	for _, flush := range flushes {
		ab.unsafeFlushAsync(ctx, flush)
	}
}

// FlushAsync clears the buffer, if a handler is provided it is passed the contents of the buffer.
// This call is asynchronous, in that it will call the flush handler on its own goroutine.
func (ab *Buffer) FlushAsync(ctx context.Context) error {
	ab.contentsMu.Lock()
	contents := ab.contents.Drain()
	ab.contentsMu.Unlock()
	ab.unsafeFlushAsync(ctx, contents)
	return nil
}

// workerAction is called by the  workers.
func (ab *Buffer) workerAction(ctx context.Context, obj interface{}) (err error) {
	typed, ok := obj.(Flush)
	if !ok {
		return fmt.Errorf("autoflush buffer; worker action argument not autoflush.Flush")
	}
	if ab.Tracer != nil {
		var finisher TraceFinisher
		ctx, finisher = ab.Tracer.StartFlush(ctx)
		defer finisher.Finish(err)
	}
	if ab.Stats != nil {
		ab.maybeStatCount(ctx, MetricFlushHandler, 1)
		start := time.Now().UTC()
		defer func() { ab.maybeStatElapsed(ctx, MetricFlushHandlerElapsed, start) }()
	}
	err = ab.Handler(typed.Context, typed.Contents)
	return
}

// returnFlushWorker returns a given worker to the worker queue.
func (ab *Buffer) returnFlushWorker(ctx context.Context, worker *async.Worker) error {
	ab.flushWorkersReady <- worker
	return nil
}

// FlushAsync clears the buffer, if a handler is provided it is passed the contents of the buffer.
// This call is asynchronous, in that it will call the flush handler on its own goroutine.
func (ab *Buffer) unsafeFlushAsync(ctx context.Context, contents []interface{}) {
	if len(contents) == 0 {
		return
	}
	if ab.Tracer != nil {
		finisher := ab.Tracer.StartQueueFlush(ctx)
		defer finisher.Finish(nil)
	}
	if ab.Stats != nil {
		ab.maybeStatCount(ctx, MetricFlush, 1)
		ab.maybeStatGauge(ctx, MetricFlushQueueLength, float64(len(ab.flushes)))
		ab.maybeStatCount(ctx, MetricFlushItemCount, len(contents))
		start := time.Now().UTC()
		defer func() {
			ab.maybeStatElapsed(ctx, MetricFlushEnqueueElapsed, start)
		}()
	}

	logger.MaybeDebugf(ab.Log, "autoflush buffer; queue flush, queue length: %d", len(ab.flushes))
	ab.flushes <- Flush{
		Context:  ctx,
		Contents: contents,
	}
}

func (ab *Buffer) maybeStatCount(ctx context.Context, metricName string, count int) {
	if ab.Stats != nil {
		_ = ab.Stats.Count(metricName, int64(count), ab.statTags(ctx)...)
	}
}

func (ab *Buffer) maybeStatGauge(ctx context.Context, metricName string, gauge float64) {
	if ab.Stats != nil {
		_ = ab.Stats.Gauge(metricName, gauge, ab.statTags(ctx)...)
	}
}

func (ab *Buffer) maybeStatElapsed(ctx context.Context, metricName string, start time.Time) {
	if ab.Stats != nil {
		elapsed := time.Now().UTC().Sub(start.UTC())
		_ = ab.Stats.Gauge(metricName, timeutil.Milliseconds(elapsed), ab.statTags(ctx)...)
		_ = ab.Stats.TimeInMilliseconds(metricName, elapsed, ab.statTags(ctx)...)
		_ = ab.Stats.Distribution(metricName, timeutil.Milliseconds(elapsed), ab.statTags(ctx)...)
	}
}

func (ab *Buffer) statTags(ctx context.Context) (tags []string) {
	if ab.Log != nil {
		ctx = ab.Log.ApplyContext(ctx)
	}
	labels := logger.GetLabels(ctx)
	for key, value := range labels {
		tags = append(tags, stats.Tag(key, value))
	}
	return
}

// Flush is an inflight flush attempt.
type Flush struct {
	Context  context.Context
	Contents []interface{}
}

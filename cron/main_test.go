package cron

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

// TestMain is the testing entrypoint.
func TestMain(m *testing.M) {
	assert.Main(m)
}

func noop(_ context.Context) error { return nil }

const (
	runAtJobName = "runAt"
)

type runAtJob struct {
	RunAt       time.Time
	RunDelegate func(ctx context.Context) error
}

var (
	_ Schedule = (*runAt)(nil)
)

type runAt time.Time

func (ra runAt) Next(after time.Time) time.Time {
	if after.Before(time.Time(ra)) {
		return time.Time(ra)
	}
	return time.Time{}
}

func (raj *runAtJob) Name() string {
	return "runAt"
}

func (raj *runAtJob) Schedule() Schedule {
	return runAt(raj.RunAt)
}

func (raj *runAtJob) Execute(ctx context.Context) error {
	return raj.RunDelegate(ctx)
}

var (
	_ Job               = (*testJobWithTimeout)(nil)
	_ ConfigProvider    = (*testJobWithTimeout)(nil)
	_ LifecycleProvider = (*testJobWithTimeout)(nil)
)

type testJobWithTimeout struct {
	RunAt                time.Time
	TimeoutDuration      time.Duration
	RunDelegate          func(ctx context.Context) error
	CancellationDelegate func()
}

func (tj *testJobWithTimeout) Name() string {
	return "testJobWithTimeout"
}

func (tj *testJobWithTimeout) Config() JobConfig {
	return JobConfig{
		Timeout: tj.TimeoutDuration,
	}
}

func (tj *testJobWithTimeout) Schedule() Schedule {
	return Immediately()
}

func (tj *testJobWithTimeout) Execute(ctx context.Context) error {
	return tj.RunDelegate(ctx)
}

func (tj *testJobWithTimeout) Lifecycle() JobLifecycle {
	return JobLifecycle{
		OnCancellation: tj.OnCancellation,
	}
}

func (tj *testJobWithTimeout) OnCancellation(ctx context.Context) {
	tj.CancellationDelegate()
}

type testJobInterval struct {
	RunEvery    time.Duration
	RunDelegate func(ctx context.Context) error
}

func (tj *testJobInterval) Name() string {
	return "testJobInterval"
}

func (tj *testJobInterval) Schedule() Schedule {
	return Every(tj.RunEvery)
}

func (tj *testJobInterval) Execute(ctx context.Context) error {
	return tj.RunDelegate(ctx)
}

var (
	_ Job            = (*testWithDisabled)(nil)
	_ ConfigProvider = (*testWithDisabled)(nil)
)

type testWithDisabled struct {
	disabled bool
	action   func(context.Context) error
}

func (twe testWithDisabled) Name() string {
	return "testWithEnabled"
}

func (twe testWithDisabled) Schedule() Schedule {
	return nil
}

func (twe testWithDisabled) Config() JobConfig {
	return JobConfig{
		Disabled: &twe.disabled,
	}
}

func (twe testWithDisabled) Execute(ctx context.Context) error {
	if twe.action != nil {
		return twe.action(ctx)
	}
	return nil
}

type mockTracer struct {
	OnStart  func(context.Context)
	OnFinish func(context.Context, error)
}

func (mt mockTracer) Start(ctx context.Context) (context.Context, TraceFinisher) {
	if mt.OnStart != nil {
		mt.OnStart(ctx)
	}
	return ctx, &mockTraceFinisher{Parent: &mt}
}

type mockTraceFinisher struct {
	Parent *mockTracer
}

func (mtf mockTraceFinisher) Finish(ctx context.Context, err error) {
	if mtf.Parent != nil && mtf.Parent.OnFinish != nil {
		mtf.Parent.OnFinish(ctx, err)
	}
}

var (
	_ Job               = (*lifecycleTest)(nil)
	_ LifecycleProvider = (*lifecycleTest)(nil)
)

func newLifecycleTest(action func(context.Context) error) *lifecycleTest {
	return &lifecycleTest{
		CompleteSignal: make(chan struct{}),
		BrokenSignal:   make(chan struct{}),
		FixedSignal:    make(chan struct{}),
		Action:         action,
	}
}

type lifecycleTest struct {
	sync.Mutex
	Starts         int
	Completes      int
	Failures       int
	CompleteSignal chan struct{}
	BrokenSignal   chan struct{}
	FixedSignal    chan struct{}
	Action         func(context.Context) error
}

func (job *lifecycleTest) Name() string { return "lifecycle-test" }

func (job *lifecycleTest) Execute(ctx context.Context) error {
	return job.Action(ctx)
}

func (job *lifecycleTest) Lifecycle() JobLifecycle {
	return JobLifecycle{
		OnBegin:    job.OnBegin,
		OnError:    job.OnError,
		OnComplete: job.OnComplete,
		OnBroken:   job.OnBroken,
		OnFixed:    job.OnFixed,
	}
}

func (job *lifecycleTest) OnBegin(ctx context.Context) {
	job.Starts++
}
func (job *lifecycleTest) OnError(ctx context.Context) {
	job.Failures++
}
func (job *lifecycleTest) OnComplete(ctx context.Context) {
	job.Lock()
	defer job.Unlock()
	close(job.CompleteSignal)
	job.CompleteSignal = make(chan struct{})
	job.Completes++
}
func (job *lifecycleTest) OnBroken(ctx context.Context) {
	job.Lock()
	defer job.Unlock()
	close(job.BrokenSignal)
	job.BrokenSignal = make(chan struct{})
}
func (job *lifecycleTest) OnFixed(ctx context.Context) {
	job.Lock()
	defer job.Unlock()
	close(job.FixedSignal)
	job.BrokenSignal = make(chan struct{})
}

type loadJobTestMinimum struct{}

func (job loadJobTestMinimum) Name() string                    { return "load-job-test-minimum" }
func (job loadJobTestMinimum) Execute(_ context.Context) error { return nil }

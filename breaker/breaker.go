package breaker

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/blend/go-sdk/ex"
)

type (
	// Action is a piece of code to run.
	Action func(context.Context) (interface{}, error)
	// OnStateChangeHandler is called when the state changes.
	OnStateChangeHandler func(ctx context.Context, from, to State, generation int64)
	// ShouldOpenProvider returns if the breaker should open.
	ShouldOpenProvider func(ctx context.Context, counts Counts) bool
	// NowProvider returns the current time.
	NowProvider func() time.Time
)

// MustNew returns a new breaker and panics if there is a construction error.
func MustNew(options ...Option) *Breaker {
	b, err := New(options...)
	if err != nil {
		panic(err)
	}
	return b
}

// New creates a new breaker with the given options.
func New(options ...Option) (*Breaker, error) {
	b := Breaker{
		ClosedExpiryInterval: DefaultClosedExpiryInterval,
		OpenExpiryInterval:   DefaultOpenExpiryInterval,
		HalfOpenMaxActions:   DefaultHalfOpenMaxActions,
	}
	for _, opt := range options {
		if err := opt(&b); err != nil {
			return nil, err
		}
	}
	return &b, nil
}

// Breaker is a state machine to prevent performing actions that are likely to fail.
type Breaker struct {
	sync.Mutex

	// OpenAction is an optional action to be called when the breaker is open (i.e. preventing calls
	// to the main action handler.)
	OpenAction Action

	// OnStateChange is an optional handler called when the breaker transitions state.
	OnStateChange OnStateChangeHandler
	// ShouldOpenProvider is called optionally to determine if we should open the breaker.
	ShouldOpenProvider ShouldOpenProvider
	// NowProvider lets you optionally inject the current time for testing.
	NowProvider NowProvider

	// HalfOpenMaxActions is the maximum number of requests
	// we can make when the state is HalfOpen.
	HalfOpenMaxActions int64
	// ClosedExpiryInterval is the cyclic period of the closed state for the CircuitBreaker to clear the internal Counts.
	// If Interval is 0, the CircuitBreaker doesn't clear internal Counts during the closed state.
	ClosedExpiryInterval time.Duration
	// OpenExpiryInterval is the period of the open state,
	// after which the state of the CircuitBreaker becomes half-open.
	// If Timeout is 0, the timeout value of the CircuitBreaker is set to 60 seconds.
	OpenExpiryInterval time.Duration
	// Counts are stats for the breaker.
	Counts Counts

	// state is the current Breaker state (Closed, HalfOpen, Open etc.)
	state State
	// generation is the current state generation.
	generation int64
	// stateExpiresAt is the time when the current state will expire.
	// It is set when we change state according to the interval
	// and the current time.
	stateExpiresAt time.Time
}

// Do runs the given action if the Breaker accepts it.
// Do returns an error instantly if the Breaker rejects the request.
// Otherwise, Do returns the result of the request.
// If a panic occurs in the request, the Breaker handles it as an error.
func (b *Breaker) Do(ctx context.Context, action Action) (interface{}, error) {
	generation, err := b.beforeAction(ctx)
	if err != nil {
		if b.OpenAction != nil {
			return b.OpenAction(ctx)
		}
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			b.afterAction(ctx, generation, false)
		}
	}()

	res, err := action(ctx)
	b.afterAction(ctx, generation, err == nil)
	return res, err
}

// EvaluateState returns the current state of the CircuitBreaker.
func (b *Breaker) EvaluateState(ctx context.Context) State {
	b.Lock()
	defer b.Unlock()

	now := time.Now()
	state, _ := b.evaluateState(ctx, now)
	return state
}

func (b *Breaker) beforeAction(ctx context.Context) (int64, error) {
	b.Lock()
	defer b.Unlock()

	now := b.now()
	state, generation := b.evaluateState(ctx, now)

	if state == StateOpen {
		return generation, ex.New(ErrOpenState)
	} else if state == StateHalfOpen && b.Counts.Requests >= b.HalfOpenMaxActions {
		return generation, ex.New(ErrTooManyRequests)
	}

	atomic.AddInt64(&b.Counts.Requests, 1)
	return generation, nil
}

func (b *Breaker) afterAction(ctx context.Context, generation int64, success bool) {
	b.Lock()
	defer b.Unlock()

	now := b.now()
	state, generation := b.evaluateState(ctx, now)
	if generation != generation {
		return
	}

	if success {
		b.success(ctx, state, now)
		return
	}
	b.failure(ctx, state, now)
}

func (b *Breaker) success(ctx context.Context, state State, now time.Time) {
	switch state {
	case StateClosed:
		atomic.AddInt64(&b.Counts.TotalSuccesses, 1)
		atomic.AddInt64(&b.Counts.ConsecutiveSuccesses, 1)
		atomic.StoreInt64(&b.Counts.ConsecutiveFailures, 0)
	case StateHalfOpen:
		atomic.AddInt64(&b.Counts.TotalSuccesses, 1)
		atomic.AddInt64(&b.Counts.ConsecutiveSuccesses, 1)
		atomic.StoreInt64(&b.Counts.ConsecutiveFailures, 0)
		if b.Counts.ConsecutiveSuccesses >= b.HalfOpenMaxActions {
			b.setState(ctx, StateClosed, now)
		}
	}
}

func (b *Breaker) failure(ctx context.Context, state State, now time.Time) {
	switch state {
	case StateClosed:
		atomic.AddInt64(&b.Counts.TotalFailures, 1)
		atomic.AddInt64(&b.Counts.ConsecutiveFailures, 1)
		atomic.StoreInt64(&b.Counts.ConsecutiveSuccesses, 0)
		if b.shouldOpen(ctx) {
			b.setState(ctx, StateOpen, now)
		}
	case StateHalfOpen:
		b.setState(ctx, StateOpen, now)
	}
}

func (b *Breaker) evaluateState(ctx context.Context, t time.Time) (state State, generation int64) {
	switch b.state {
	case StateClosed:
		if !b.stateExpiresAt.IsZero() && b.stateExpiresAt.Before(t) {
			b.incrementGeneration(t)
		}
	case StateOpen:
		if b.stateExpiresAt.Before(t) {
			b.setState(ctx, StateHalfOpen, t)
		}
	}
	return b.state, b.generation
}

func (b *Breaker) setState(ctx context.Context, state State, now time.Time) {
	if b.state == state {
		return
	}

	previousState := b.state
	b.state = state
	b.incrementGeneration(now)
	if b.OnStateChange != nil {
		b.OnStateChange(ctx, previousState, b.state, b.generation)
	}
}

func (b *Breaker) incrementGeneration(now time.Time) {
	atomic.AddInt64(&b.generation, 1)
	b.Counts = Counts{}

	var zero time.Time
	switch b.state {
	case StateClosed:
		if b.ClosedExpiryInterval == 0 {
			b.stateExpiresAt = zero
		} else {
			b.stateExpiresAt = now.Add(b.ClosedExpiryInterval)
		}
	case StateOpen:
		b.stateExpiresAt = now.Add(b.OpenExpiryInterval)
	default: // StateHalfOpen
		b.stateExpiresAt = zero
	}
}

func (b *Breaker) shouldOpen(ctx context.Context) bool {
	if b.ShouldOpenProvider != nil {
		return b.ShouldOpenProvider(ctx, b.Counts)
	}
	return b.Counts.ConsecutiveFailures > DefaultConsecutiveFailures
}

func (b *Breaker) now() time.Time {
	if b.NowProvider != nil {
		return b.NowProvider()
	}
	return time.Now()
}

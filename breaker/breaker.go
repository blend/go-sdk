package breaker

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/blend/go-sdk/ex"
)

// New creates a new breaker.
func New(action func(context.Context) error, options ...Option) (*Breaker, error) {
	b := Breaker{
		Action:               action,
		ClosedExpiryInterval: DefaultClosedExpiryInterval,
		OpenExpiryInterval:   DefaultOpenExpiryInterval,
		HalfOpenMaxRequests:  DefaultHalfOpenMaxRequests,
	}
	for _, opt := range options {
		if err := opt(&b); err != nil {
			return nil, err
		}
	}
	return &b, nil
}

// Breaker is a state machine to prevent sending requests that are likely to fail.
type Breaker struct {
	sync.Mutex

	Action              func(context.Context) error
	OnStateChange       func(ctx context.Context, from, to State, generation int64)
	ShouldCloseProvider func(ctx context.Context, counts Counts) bool
	NowProvider         func() time.Time

	// HalfOpenMaxRequests is the maximum number of requests
	// we can make when the state is HalfOpen.
	HalfOpenMaxRequests int64

	// ClosedExpiryInterval is the cyclic period of the closed state for the CircuitBreaker to clear the internal Counts.
	// If Interval is 0, the CircuitBreaker doesn't clear internal Counts during the closed state.
	ClosedExpiryInterval time.Duration

	// OpenExpiryInterval is the period of the open state,
	// after which the state of the CircuitBreaker becomes half-open.
	// If Timeout is 0, the timeout value of the CircuitBreaker is set to 60 seconds.
	OpenExpiryInterval time.Duration

	// State is the current breaker state (Open, HalfOpen, Closed).
	// Counts are stats for the breaker.
	Counts Counts

	// State is the current Breaker state (Closed, HalfOpen, Open etc.)
	State State

	// Generation is the current state generation.
	Generation int64

	// StateExpiresAt is the time when the current state will expire.
	// It is set when we change state according to the interval
	// and the current time.
	StateExpiresAt time.Time
}

// Execute runs the given request if the CircuitBreaker accepts it.
// Execute returns an error instantly if the CircuitBreaker rejects the request.
// Otherwise, Execute returns the result of the request.
// If a panic occurs in the request, the CircuitBreaker handles it as an error
// and causes the same panic again.
func (b *Breaker) Execute(ctx context.Context) error {
	generation, err := b.beforeAction(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			b.afterAction(ctx, generation, false)
		}
	}()

	err = b.Action(ctx)
	b.afterAction(ctx, generation, err == nil)
	return err
}

func (b *Breaker) beforeAction(ctx context.Context) (int64, error) {
	b.Lock()
	defer b.Unlock()

	now := b.now()
	ctx = WithContextTime(ctx, now)
	state, generation := b.getState(ctx, now)

	if state == StateOpen {
		return generation, ex.New(ErrOpenState)
	} else if state == StateHalfOpen && b.Counts.Requests >= b.HalfOpenMaxRequests {
		return generation, ex.New(ErrTooManyRequests)
	}

	atomic.AddInt64(&b.Counts.Requests, 1)
	return generation, nil
}

func (b *Breaker) afterAction(ctx context.Context, generation int64, success bool) {
	b.Lock()
	defer b.Unlock()

	now := b.now()
	ctx = WithContextTime(ctx, now)
	state, generation := b.getState(ctx, now)
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
		if b.Counts.ConsecutiveSuccesses >= b.HalfOpenMaxRequests {
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
		if b.shouldClose(ctx) {
			b.setState(ctx, StateOpen, now)
		}
	case StateHalfOpen:
		b.setState(ctx, StateOpen, now)
	}
}

func (b *Breaker) getState(ctx context.Context, t time.Time) (state State, generation int64) {
	switch b.State {
	case StateClosed:
		if !b.StateExpiresAt.IsZero() && b.StateExpiresAt.Before(t) {
			b.incrementGeneration(t)
		}
	case StateOpen:
		if b.StateExpiresAt.Before(t) {
			b.setState(ctx, StateHalfOpen, t)
		}
	}
	return b.State, b.Generation
}

func (b *Breaker) setState(ctx context.Context, state State, now time.Time) {
	if b.State == state {
		return
	}

	previousState := b.State
	b.State = state
	b.incrementGeneration(now)
	if b.OnStateChange != nil {
		b.OnStateChange(ctx, previousState, b.State, b.Generation)
	}
}

func (b *Breaker) incrementGeneration(now time.Time) {
	atomic.AddInt64(&b.Generation, 1)

	var zero time.Time
	switch b.State {
	case StateClosed:
		if b.ClosedExpiryInterval == 0 {
			b.StateExpiresAt = zero
		} else {
			b.StateExpiresAt = now.Add(b.ClosedExpiryInterval)
		}
	case StateOpen:
		b.StateExpiresAt = now.Add(b.OpenExpiryInterval)
	default: // StateHalfOpen
		b.StateExpiresAt = zero
	}
}

func (b *Breaker) shouldClose(ctx context.Context) bool {
	if b.ShouldCloseProvider != nil {
		return b.ShouldCloseProvider(ctx, b.Counts)
	}
	return b.Counts.ConsecutiveFailures > 5
}

func (b *Breaker) now() time.Time {
	if b.NowProvider != nil {
		return b.NowProvider()
	}
	return time.Now()
}

package breaker

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	b, err := New()
	assert.Nil(err)
	assert.Equal(DefaultHalfOpenMaxActions, b.HalfOpenMaxActions)
	assert.Equal(DefaultOpenExpiryInterval, b.OpenExpiryInterval)
	assert.Equal(DefaultClosedExpiryInterval, b.ClosedExpiryInterval)
}

func TestNewOptions(t *testing.T) {
	assert := assert.New(t)

	b, err := New(
		OptHalfOpenMaxActions(5),
		OptOpenExpiryInterval(10*time.Second),
		OptClosedExpiryInterval(20*time.Second),
	)
	assert.Nil(err)
	assert.Equal(5, b.HalfOpenMaxActions)
	assert.Equal(10*time.Second, b.OpenExpiryInterval)
	assert.Equal(20*time.Second, b.ClosedExpiryInterval)
}

func createTestBreaker() *Breaker {
	return MustNew(nil, OptClosedExpiryInterval(0))
}

func succeed(ctx context.Context, b *Breaker) error {
	_, err := b.Do(ctx, func(_ context.Context) (interface{}, error) { return nil, nil })
	return err
}

func pseudoSleep(b *Breaker, period time.Duration) {
	if !b.stateExpiresAt.IsZero() {
		b.stateExpiresAt = b.stateExpiresAt.Add(-period)
	}
}

func fail(ctx context.Context, b *Breaker) error {
	msg := "fail"
	_, err := b.Do(ctx, func(_ context.Context) (interface{}, error) { return nil, fmt.Errorf(msg) })
	if err.Error() == msg {
		return nil
	}
	return err
}

func TestBreaker(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	b := createTestBreaker()

	for i := 0; i < 5; i++ {
		assert.Nil(fail(ctx, b))
	}
	assert.Equal(StateClosed, b.State(ctx))
	assert.Equal(Counts{5, 0, 5, 0, 5}, b.Counts)

	assert.Nil(succeed(ctx, b))
	assert.Equal(StateClosed, b.State(ctx))
	assert.Equal(Counts{6, 1, 5, 1, 0}, b.Counts)

	assert.Nil(fail(ctx, b))
	assert.Equal(StateClosed, b.State(ctx))
	assert.Equal(Counts{7, 1, 6, 0, 1}, b.Counts)

	// StateClosed to StateOpen
	for i := 0; i < 5; i++ {
		assert.Nil(fail(ctx, b)) // 6 consecutive failures
	}
	assert.Equal(StateOpen, b.State(ctx))
	assert.Equal(Counts{0, 0, 0, 0, 0}, b.Counts)
	assert.False(b.stateExpiresAt.IsZero())

	assert.NotNil(succeed(ctx, b))
	assert.NotNil(fail(ctx, b))
	assert.Equal(Counts{0, 0, 0, 0, 0}, b.Counts)

	pseudoSleep(b, time.Duration(59)*time.Second)
	assert.Equal(StateOpen, b.State(ctx))

	// StateOpen to StateHalfOpen
	pseudoSleep(b, time.Duration(1)*time.Second) // over Timeout
	assert.Equal(StateHalfOpen, b.State(ctx))
	assert.True(b.stateExpiresAt.IsZero())

	// StateHalfOpen to StateOpen
	assert.Nil(fail(ctx, b))
	assert.Equal(StateOpen, b.State(ctx))
	assert.Equal(Counts{0, 0, 0, 0, 0}, b.Counts)
	assert.False(b.stateExpiresAt.IsZero())

	// StateOpen to StateHalfOpen
	pseudoSleep(b, time.Duration(60)*time.Second)
	assert.Equal(StateHalfOpen, b.State(ctx))
	assert.True(b.stateExpiresAt.IsZero())

	// StateHalfOpen to StateClosed
	assert.Nil(succeed(ctx, b))
	assert.Equal(StateClosed, b.State(ctx))
	assert.Equal(Counts{0, 0, 0, 0, 0}, b.Counts)
	assert.True(b.stateExpiresAt.IsZero())
}

func TestBreakerErrStateOpen(t *testing.T) {
	assert := assert.New(t)

	var didCall bool
	b, err := New()
	assert.Nil(err)

	b.state = StateOpen
	b.stateExpiresAt = time.Now().Add(time.Hour)

	_, err = b.Do(context.Background(), func(_ context.Context) (interface{}, error) {
		didCall = true
		return nil, nil
	})
	assert.True(ex.Is(err, ErrOpenState), fmt.Sprintf("%v", err))
	assert.False(didCall)
}

func TestBreakerErrTooManyRequests(t *testing.T) {
	assert := assert.New(t)

	var didCall bool
	b, err := New()
	assert.Nil(err)

	b.state = StateHalfOpen
	b.Counts.Requests = 10
	b.HalfOpenMaxActions = 5

	_, err = b.Do(context.Background(), func(_ context.Context) (interface{}, error) {
		didCall = true
		return nil, nil
	})
	assert.True(ex.Is(err, ErrTooManyRequests))
	assert.False(didCall)
}

func TestBreakerCallsClosedHandler(t *testing.T) {
	assert := assert.New(t)

	var didCall, didCallClosed bool
	b, err := New(OptOpenAction(func(_ context.Context) {
		didCallClosed = true
	}))
	assert.Nil(err)

	b.state = StateOpen
	b.stateExpiresAt = time.Now().Add(time.Hour)

	_, err = b.Do(context.Background(), func(_ context.Context) (interface{}, error) {
		didCall = true
		return nil, nil
	})

	assert.True(ex.Is(err, ErrOpenState))
	assert.False(didCall)
	assert.True(didCallClosed)
}

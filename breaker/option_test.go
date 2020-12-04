package breaker

import (
	"context"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestOptHalfOpenMaxActions(t *testing.T) {
	assert := assert.New(t)

	b := new(Breaker)
	assert.Zero(b.HalfOpenMaxActions)
	assert.Nil(OptHalfOpenMaxActions(5)(b))
	assert.Equal(5, b.HalfOpenMaxActions)
}

func TestOptClosedExpiryInterval(t *testing.T) {
	assert := assert.New(t)

	b := new(Breaker)
	assert.Zero(b.ClosedExpiryInterval)
	assert.Nil(OptClosedExpiryInterval(5 * time.Second)(b))
	assert.Equal(5*time.Second, b.ClosedExpiryInterval)
}

func TestOptOpenExpiryInterval(t *testing.T) {
	assert := assert.New(t)

	b := new(Breaker)
	assert.Zero(b.OpenExpiryInterval)
	assert.Nil(OptOpenExpiryInterval(5 * time.Second)(b))
	assert.Equal(5*time.Second, b.OpenExpiryInterval)
}

func TestOptConfig(t *testing.T) {
	assert := assert.New(t)

	b := new(Breaker)
	assert.Zero(b.HalfOpenMaxActions)
	assert.Zero(b.ClosedExpiryInterval)
	assert.Zero(b.OpenExpiryInterval)
	assert.Nil(OptConfig(Config{
		HalfOpenMaxActions:   1,
		ClosedExpiryInterval: 2 * time.Second,
		OpenExpiryInterval:   3 * time.Second,
	})(b))
	assert.Equal(1, b.HalfOpenMaxActions)
	assert.Equal(2*time.Second, b.ClosedExpiryInterval)
	assert.Equal(3*time.Second, b.OpenExpiryInterval)
}

func TestOptOpenAction(t *testing.T) {
	assert := assert.New(t)

	b := new(Breaker)
	assert.Nil(b.OpenAction)
	assert.Nil(OptOpenAction(func(_ context.Context) (interface{}, error) { return nil, nil })(b))
	assert.NotNil(b.OpenAction)
}

func TestOptOnStateChange(t *testing.T) {
	assert := assert.New(t)

	b := new(Breaker)
	assert.Nil(b.OnStateChange)
	assert.Nil(OptOnStateChange(func(_ context.Context, from, to State, generation int64) {})(b))
	assert.NotNil(b.OnStateChange)
}

func TestOptShouldOpenProvider(t *testing.T) {
	assert := assert.New(t)

	b := new(Breaker)
	assert.Nil(b.ShouldOpenProvider)
	assert.Nil(OptShouldOpenProvider(func(ctx context.Context, counts Counts) bool { return false })(b))
	assert.NotNil(b.ShouldOpenProvider)
}

func TestOptNowProvider(t *testing.T) {
	assert := assert.New(t)

	b := new(Breaker)
	assert.Nil(b.NowProvider)
	assert.Nil(OptNowProvider(time.Now)(b))
	assert.NotNil(b.NowProvider)
}

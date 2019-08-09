package breaker

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestBreakerErrStateOpen(t *testing.T) {
	assert := assert.New(t)

	var didCall bool
	b, err := New(func(_ context.Context) error {
		return nil
	})
	assert.Nil(err)

	b.State = StateOpen
	b.StateExpiresAt = time.Now().Add(time.Hour)

	err = b.Execute(context.Background())
	assert.True(ex.Is(err, ErrOpenState), fmt.Sprintf("%v", err))
	assert.False(didCall)
}

func TestBreakerErrTooManyRequests(t *testing.T) {
	assert := assert.New(t)

	var didCall bool
	b, err := New(func(_ context.Context) error {
		return nil
	})
	assert.Nil(err)

	b.State = StateHalfOpen
	b.Counts.Requests = 10
	b.HalfOpenMaxRequests = 5

	err = b.Execute(context.Background())
	assert.True(ex.Is(err, ErrTooManyRequests))
	assert.False(didCall)
}

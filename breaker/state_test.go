package breaker

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestStateConstants(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(State(0), StateClosed)
	assert.Equal(State(1), StateHalfOpen)
	assert.Equal(State(2), StateOpen)

	assert.Equal(StateClosed.String(), "closed")
	assert.Equal(StateHalfOpen.String(), "half-open")
	assert.Equal(StateOpen.String(), "open")
	assert.Equal(State(100).String(), "unknown state: 100")
}

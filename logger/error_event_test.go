package logger

import (
	"context"
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNewErrorEvent(t *testing.T) {
	assert := assert.New(t)

	/// stuff
	ee := NewErrorEvent(Fatal, fmt.Errorf("only a test"), OptErrorEventState("foo"))
	assert.Equal(Fatal, ee.GetFlag())
	assert.Equal("only a test", ee.Err.Error())
	assert.Equal("foo", ee.State)
}

func TestErrorEventListener(t *testing.T) {
	assert := assert.New(t)

	ee := NewErrorEvent(Fatal, fmt.Errorf("only a test"))

	var didCall bool
	ml := NewErrorEventListener(func(ctx context.Context, e *ErrorEvent) {
		didCall = true
	})

	ml(context.Background(), ee)
	assert.True(didCall)
}

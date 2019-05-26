package logger

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestHTTPResponseEventListener(t *testing.T) {
	assert := assert.New(t)

	var didCall bool
	listener := NewHTTPResponseEventListener(func(_ context.Context, hre *HTTPResponseEvent) {
		didCall = true
	})
	listener(context.Background(), NewMessageEvent(Info, "test"))
	assert.False(didCall)
	listener(context.Background(), NewHTTPResponseEvent(nil))
	assert.True(didCall)
}

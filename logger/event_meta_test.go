package logger

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNewEventMeta(t *testing.T) {
	assert := assert.New(t)

	em := NewEventMeta(Info)
	assert.Equal(Info, em.GetFlag())
	assert.False(em.GetTimestamp().IsZero())
}

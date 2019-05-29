package logger

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestMarshalEvent(t *testing.T) {
	assert := assert.New(t)

	typed, ok := MarshalEvent("foo")
	assert.Nil(typed)
	assert.False(ok)

	typed, ok = MarshalEvent(NewMessageEvent(Info, "hi"))
	assert.True(ok)
	assert.NotNil(typed)
}

package configutil

import (
	"context"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestDurationPtr(t *testing.T) {
	assert := assert.New(t)

	isNil := DurationPtr(nil)
	value := time.Second
	hasValue := DurationPtr(&value)
	value2 := time.Millisecond
	hasValue2 := DurationPtr(&value2)

	var setValue time.Duration
	assert.Nil(SetDuration(&setValue, isNil, hasValue, hasValue2)(context.TODO()))
	assert.Equal(time.Second, setValue)
}

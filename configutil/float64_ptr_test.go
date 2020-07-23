package configutil

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestFloat64Ptr(t *testing.T) {
	assert := assert.New(t)

	isNil := Float64Ptr(nil)
	value := 1.0
	hasValue := Float64Ptr(&value)
	value2 := 2.0
	hasValue2 := Float64Ptr(&value2)

	var setValue float64
	assert.Nil(SetFloat64(&setValue, isNil, hasValue, hasValue2)(context.TODO()))
	assert.Equal(1.0, setValue)
}

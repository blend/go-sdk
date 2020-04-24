package configutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestIntPtr(t *testing.T) {
	assert := assert.New(t)

	isNil := IntPtr(nil)
	value := 1
	hasValue := IntPtr(&value)
	value2 := 2
	hasValue2 := IntPtr(&value2)

	var setValue int
	assert.Nil(SetInt(&setValue, isNil, hasValue, hasValue2))
	assert.Equal(1, setValue)
}

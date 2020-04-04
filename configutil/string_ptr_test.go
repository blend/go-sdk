package configutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestStringPtr(t *testing.T) {
	assert := assert.New(t)

	isNil := StringPtr(nil)
	emptyValue := ""
	isEmpty := StringPtr(&emptyValue)
	value := "foo"
	hasValue := StringPtr(&value)
	value2 := "bar"
	hasValue2 := StringPtr(&value2)

	var setValue string
	assert.Nil(SetString(&setValue, isNil, isEmpty, hasValue, hasValue2))
	assert.Equal("foo", setValue)
}

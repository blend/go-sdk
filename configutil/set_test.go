package configutil

import (
	"testing"
	"time"
	"github.com/blend/go-sdk/assert"
)

func TestSetString(t *testing.T) {
	assert := assert.New(t)

	empty := String("")
	hasValue := String("has value")
	hasValue2 := String("has another value")

	var value string
	assert.Nil(SetString(&value, empty, hasValue, hasValue2))
	assert.Equal("has value", value)
}

func TestSetStrings(t *testing.T) {
	assert := assert.New(t)

	empty := Strings(nil)
	hasValue := Strings([]string{"has value"})
	hasValue2 := Strings([]string{"has another value"})

	var value []string
	assert.Nil(SetStrings(&value, empty, hasValue, hasValue2))
	assert.Equal([]string{"has value"}, value)
}

func TestSetBool(t *testing.T) {
	assert := assert.New(t)

	empty := Bool(nil)
	tv := true
	hasValue := Bool(&tv)
	fv := false
	hasValue2 := Bool(&fv)

	var value *bool
	assert.Nil(SetBool(&value, empty, hasValue, hasValue2))
	assert.Equal(true, *value)
}

func TestSetInt(t *testing.T) {
	assert := assert.New(t)

	empty := Parse(String(""))
	hasValue := Int(1)
	hasValue2 := Int(2)

	var value int
	assert.Nil(SetInt(&value, empty, hasValue, hasValue2))
	assert.Equal(1, value)

	errors := Parse(String("bad"))
	assert.NotNil(SetInt(&value, errors))
}

func TestSetFloat64(t *testing.T) {
	assert := assert.New(t)

	empty := Parse(String(""))
	hasValue := Float64(1)
	hasValue2 := Float64(2)

	var value float64
	assert.Nil(SetFloat64(&value, empty, hasValue, hasValue2))
	assert.Equal(1, value)

	errors := Parse(String("bad"))
	assert.NotNil(SetFloat64(&value, errors))
}

func TestSetDuration(t *testing.T) {
	assert := assert.New(t)

	empty := Parse(String(""))
	hasValue := Duration(time.Second)
	hasValue2 := Duration(2 * time.Second)

	var value time.Duration
	assert.Nil(SetDuration(&value, empty, hasValue, hasValue2))
	assert.Equal(time.Second, value)

	errors := Parse(String("bad"))
	assert.NotNil(SetDuration(&value, errors))
}

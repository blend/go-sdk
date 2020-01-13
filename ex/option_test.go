package ex

import (
	"fmt"
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestOptMessage(t *testing.T) {
	assert := assert.New(t)

	ex := &Ex{}

	OptMessage("a message", " bar")(ex)
	assert.Equal("a message bar", ex.Message)
}

func TestOptMessagef(t *testing.T) {
	assert := assert.New(t)

	ex := &Ex{}

	OptMessagef("a message %s", "bar")(ex)
	assert.Equal("a message bar", ex.Message)
}

func TestOptStackTrace(t *testing.T) {
	assert := assert.New(t)

	ex := &Ex{}

	OptStackTrace(StackStrings([]string{"first", "second"}))(ex)
	assert.NotNil(ex.StackTrace)
	assert.Equal([]string{"first", "second"}, ex.StackTrace.Strings())
}

func TestOptInner(t *testing.T) {
	assert := assert.New(t)

	ex := &Ex{}

	OptInner(fmt.Errorf("this is only a test"))(ex)
	assert.NotNil(ex.Inner)
}

func TestOptInnerClass(t *testing.T) {
	assert := assert.New(t)

	ex := &Ex{}

	OptInnerClass(fmt.Errorf("this is only a test"))(ex)
	assert.NotNil(ex.Inner)
	assert.Nil(ErrStackTrace(ex.Inner))
}

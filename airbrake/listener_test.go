package airbrake

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestFrames(t *testing.T) {
	assert := assert.New(t)

	exErr := ex.As(ex.New("this is a test"))
	fr := frames(exErr.Stack)
	assert.NotEmpty(fr, fmt.Sprintf("%T", exErr.Stack))

	fr = frames(ex.As(ex.New("this is a test", ex.OptStack(ex.StackStrings([]string{"foo", "bar"})))).Stack)
	assert.Empty(fr)
}

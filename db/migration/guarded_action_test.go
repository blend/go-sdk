package migration

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestStep(t *testing.T) {
	assert := assert.New(t)

	step := NewStep(Always(), NoOp)
	assert.NotNil(step.Guard)
	assert.NotNil(step.Body)
}
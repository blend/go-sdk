package migration

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNewGroup(t *testing.T) {
	assert := assert.New(t)

	g := NewGroup(OptGroupActions(NewStep(Always(), NoOp)))
	assert.Len(g.Actions, 1)
}

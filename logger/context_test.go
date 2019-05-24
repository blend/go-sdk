package logger

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestContext(t *testing.T) {
	assert := assert.New(t)

	log := None()
	ctx := NewContext(log, []string{"foo", "bar"}, Fields{"zoo": "who"}, OptContextSetPath("bar", "bazz"), OptContextFields(Fields{"moo": "loo"}))
	assert.NotNil(ctx.Logger)
	assert.Equal([]string{"bar", "bazz"}, ctx.Path)
	assert.Equal("loo", ctx.Fields["moo"])

	subCtx := ctx.SubContext("bailey").WithFields(Fields{"what": "where"})
	assert.Equal([]string{"bar", "bazz", "bailey"}, subCtx.Path)
	assert.Equal("where", subCtx.Fields["what"])
	assert.Equal("loo", subCtx.Fields["moo"])
}

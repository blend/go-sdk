package logger

import (
	"context"
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

func TestContextTrigger(t *testing.T) {
	assert := assert.New(t)

	log := MustNew(OptEnabled("test"))
	log.Output = nil
	fired := make(chan struct{})
	var contextPath []string
	var contextFields Fields
	log.Listen("test", DefaultListenerName, func(ctx context.Context, e Event) {
		defer close(fired)
		contextPath, contextFields = GetSubContextMeta(ctx)
	})
	ctx := NewContext(log, []string{"path"}, Fields{"one": "two"})

	ctx.Trigger(context.Background(), NewMessageEvent("test", "this is only a test"))
	<-fired

	assert.Equal([]string{"path"}, contextPath)
	assert.Equal(Fields{"one": "two"}, contextFields)
}

func TestContextSyncTrigger(t *testing.T) {
	assert := assert.New(t)

	log := MustNew(OptEnabled("test"))
	log.Output = nil
	fired := make(chan struct{})
	var contextPath []string
	var contextFields Fields
	log.Listen("test", DefaultListenerName, func(ctx context.Context, e Event) {
		defer close(fired)
		contextPath, contextFields = GetSubContextMeta(ctx)
	})
	ctx := NewContext(log, []string{"path"}, Fields{"one": "two"})

	go ctx.SyncTrigger(context.Background(), NewMessageEvent("test", "this is only a test"))
	<-fired

	assert.Equal([]string{"path"}, contextPath)
	assert.Equal(Fields{"one": "two"}, contextFields)
}

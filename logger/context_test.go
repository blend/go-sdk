package logger

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
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

func TestOptContextPath(t *testing.T) {
	assert := assert.New(t)

	log := None()
	sc := log.SubContext("foo", OptContextPath("bar"))
	assert.Equal([]string{"foo", "bar"}, sc.Path)
}

func TestOptContextSetFields(t *testing.T) {
	assert := assert.New(t)

	log := None()
	log.Fields = Fields{"foo": "far"}
	sc := log.SubContext("path", OptContextSetFields(Fields{"foo": "bar"}))
	assert.Equal("bar", sc.Fields["foo"])
}

func TestContextMethods(t *testing.T) {
	assert := assert.New(t)

	log := All()
	log.Formatter = NewTextOutputFormatter(OptTextNoColor(), OptTextHideTimestamp())

	buf := new(bytes.Buffer)
	log.Output = buf
	log.Info("format", " test")
	assert.Equal("[info] format test\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	log.Debug("format", " test")
	assert.Equal("[debug] format test\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	log.WarningWithReq(fmt.Errorf("only a test"), &http.Request{Method: "foo"})
	assert.Equal("[warning] only a test\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	log.ErrorWithReq(fmt.Errorf("only a test"), &http.Request{Method: "foo"})
	assert.Equal("[error] only a test\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	log.FatalWithReq(fmt.Errorf("only a test"), &http.Request{Method: "foo"})
	assert.Equal("[fatal] only a test\n", buf.String())
}

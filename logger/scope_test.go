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
	ctx := NewScope(log, []string{"foo", "bar"}, Labels{"zoo": "who"}, Annotations{"annotation0": "one"}, OptScopeSetPath("bar", "bazz"), OptScopeLabels(Labels{"moo": "loo"}))
	assert.NotNil(ctx.Logger)
	assert.Equal([]string{"bar", "bazz"}, ctx.Path)
	assert.Equal("loo", ctx.Labels["moo"])

	subCtx := ctx.SubScope("bailey").WithLabels(Labels{"what": "where"})
	assert.Equal([]string{"bar", "bazz", "bailey"}, subCtx.Path)
	assert.Equal("where", subCtx.Labels["what"])
	assert.Equal("loo", subCtx.Labels["moo"])
}

func TestContextTrigger(t *testing.T) {
	assert := assert.New(t)

	log := MustNew(OptEnabled("test"))
	log.Output = nil
	fired := make(chan struct{})
	var scopePath []string
	var scopeLabels Labels
	var scopeAnnotations Annotations
	log.Listen("test", DefaultListenerName, func(ctx context.Context, e Event) {
		defer close(fired)
		scopePath, scopeLabels, scopeAnnotations = GetSubScopeMeta(ctx)
	})
	ctx := NewScope(log, []string{"path"}, Labels{"one": "two"}, Annotations{"three": "four"})

	ctx.Trigger(context.Background(), NewMessageEvent("test", "this is only a test"))
	<-fired

	assert.Equal([]string{"path"}, scopePath)
	assert.Equal(Labels{"one": "two"}, scopeLabels)
	assert.Equal(Annotations{"three": "four"}, scopeAnnotations)
}

func TestContextSyncTrigger(t *testing.T) {
	assert := assert.New(t)

	log := MustNew(OptEnabled("test"))
	log.Output = nil
	fired := make(chan struct{})
	var scopePath []string
	var scopeLabels Labels
	var scopeAnnotations Annotations
	log.Listen("test", DefaultListenerName, func(ctx context.Context, e Event) {
		defer close(fired)
		scopePath, scopeLabels, scopeAnnotations = GetSubScopeMeta(ctx)
	})
	ctx := NewScope(log, []string{"path"}, Labels{"one": "two"}, Annotations{"three": "four"})

	go ctx.SyncTrigger(context.Background(), NewMessageEvent("test", "this is only a test"))
	<-fired

	assert.Equal([]string{"path"}, scopePath)
	assert.Equal(Labels{"one": "two"}, scopeLabels)
	assert.Equal(Annotations{"three": "four"}, scopeAnnotations)
}

func TestOptContextPath(t *testing.T) {
	assert := assert.New(t)

	log := None()
	sc := log.SubScope("foo", OptScopePath("bar"))
	assert.Equal([]string{"foo", "bar"}, sc.Path)
}

func TestOptContextSetFields(t *testing.T) {
	assert := assert.New(t)

	log := None()
	log.Labels = Labels{"foo": "far"}
	sc := log.SubScope("path", OptScopeSetLabels(Labels{"foo": "bar"}))
	assert.Equal("bar", sc.Labels["foo"])
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

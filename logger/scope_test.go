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
	ctx := NewScope(
		WithLabels(WithScopePath(context.Background(), "foo", "bar"), Labels{"moo": "loo"}),
		log,
	)
	assert.NotNil(ctx.Logger)
	assert.Equal([]string{"foo", "bar"}, GetScopePath(ctx.Context()))
	assert.Equal("loo", GetLabels(ctx.Context())["moo"])

	subCtx := ctx.WithPath("bailey").WithLabels(Labels{"what": "where"})
	assert.Equal([]string{"foo", "bar", "bailey"}, GetScopePath(subCtx.Context()))
	assert.Equal("where", GetLabels(subCtx.Context())["what"])
	assert.Equal("loo", GetLabels(subCtx.Context())["moo"])
}

func TestWithPath(t *testing.T) {
	assert := assert.New(t)

	log := None()
	sc := log.WithPath("foo", "bar")
	assert.Equal([]string{"foo", "bar"}, GetScopePath(sc.Context()))
}

func TestWithLabels(t *testing.T) {
	assert := assert.New(t)

	log := None()
	sc := log.WithLabels(Labels{"foo": "bar"})
	assert.Equal("bar", GetLabels(sc.Context())["foo"])
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

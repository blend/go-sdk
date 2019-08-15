package logger

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNewScope(t *testing.T) {
	assert := assert.New(t)

	log := None()
	ctx := NewScope(
		WithFields(WithScopePath(context.Background(), "foo", "bar"), Fields{"moo": "loo"}),
		log,
	)
	assert.NotNil(ctx.Logger)
	assert.Equal([]string{"foo", "bar"}, GetScopePath(ctx.Context()))
	assert.Equal("loo", GetFields(ctx.Context())["moo"])

	subCtx := ctx.WithPath("bailey").WithFields(Fields{"what": "where"})
	assert.Equal([]string{"foo", "bar", "bailey"}, GetScopePath(subCtx.Context()))
	assert.Equal("where", GetFields(subCtx.Context())["what"])
	assert.Equal("loo", GetFields(subCtx.Context())["moo"])
}

func TestWithPath(t *testing.T) {
	assert := assert.New(t)

	log := None()
	sc := log.WithPath("foo", "bar")
	assert.Equal([]string{"foo", "bar"}, GetScopePath(sc.Context()))
}

func TestWithFields(t *testing.T) {
	assert := assert.New(t)

	log := None()
	sc := log.WithFields(Fields{"foo": "bar"})
	assert.Equal("bar", GetFields(sc.Context())["foo"])
}

func TestScopeMethods(t *testing.T) {
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

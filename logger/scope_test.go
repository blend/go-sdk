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
	sc := NewScope(
		log,
		OptScopePath("foo", "bar"),
		OptScopeFields(Fields{"moo": "loo"}),
	)
	assert.NotNil(sc.Logger)
	assert.Equal([]string{"foo", "bar"}, sc.Path)
	assert.Equal("loo", sc.Fields["moo"])

	sub := sc.WithPath("bailey").WithFields(Fields{"what": "where"})
	assert.Equal([]string{"foo", "bar", "bailey"}, sub.Path)
	assert.Equal("where", sub.Fields["what"])
	assert.Equal("loo", sub.Fields["moo"])
}

func TestWithPath(t *testing.T) {
	assert := assert.New(t)

	log := None()
	sc := log.WithPath("foo", "bar")
	assert.Equal([]string{"foo", "bar"}, sc.Path)
}

func TestWithFields(t *testing.T) {
	assert := assert.New(t)

	log := None()
	sc := log.WithFields(Fields{"foo": "bar"})
	assert.Equal("bar", sc.Fields["foo"])
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

	buf = new(bytes.Buffer)
	log.Output = buf
	log.Path = []string{"outer", "inner"}
	log.Fields = Fields{"foo": "bar"}
	log.Info("format test")
	assert.Equal("[outer > inner] [info] format test\tfoo=bar\n", buf.String())

}

func TestScopeApplyContext(t *testing.T) {
	assert := assert.New(t)

	sc := NewScope(None())
	sc.Path = []string{"one", "two"}
	sc.Fields = Fields{"foo": "bar"}

	ctx := WithFields(context.Background(), Fields{"moo": "loo"})
	ctx = WithScopePath(ctx, "three", "four")

	final := sc.ApplyContext(ctx)
	assert.Equal([]string{"one", "two", "three", "four"}, GetScopePath(final))
	assert.Equal("bar", GetFields(final)["foo"])
	assert.Equal("loo", GetFields(final)["moo"])
}

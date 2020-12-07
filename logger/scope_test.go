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
		OptScopeLabels(Labels{"moo": "loo"}),
		OptScopeAnnotations(Annotations{"alpha": "bravo"}),
	)
	assert.NotNil(sc.Logger)
	assert.Equal([]string{"foo", "bar"}, sc.Path)
	assert.Equal("loo", sc.Labels["moo"])

	sub := sc.WithPath("example-string").WithLabels(Labels{"what": "where"}).WithAnnotations(Annotations{"zoo": 47})
	assert.Equal([]string{"foo", "bar", "example-string"}, sub.Path)
	assert.Equal("where", sub.Labels["what"])
	assert.Equal("loo", sub.Labels["moo"])
	assert.Equal(47, sub.Annotations["zoo"])
	assert.Equal("bravo", sub.Annotations["alpha"])
}

func TestWithPath(t *testing.T) {
	assert := assert.New(t)

	log := None()
	sc := log.WithPath("foo", "bar")
	assert.Equal([]string{"foo", "bar"}, sc.Path)
}

func TestWithLabels(t *testing.T) {
	assert := assert.New(t)

	log := None()
	sc := log.WithLabels(Labels{"foo": "bar"})
	assert.Equal("bar", sc.Labels["foo"])
}

func TestWithAnnotations(t *testing.T) {
	assert := assert.New(t)

	log := None()
	sc := log.WithAnnotations(Annotations{"foo": "bar"})
	assert.Equal("bar", sc.Annotations["foo"])
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
	log.Warning(fmt.Errorf("only a test"), OptErrorEventState(&http.Request{Method: "foo"}))
	assert.Equal("[warning] only a test\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	log.Error(fmt.Errorf("only a test"), OptErrorEventState(&http.Request{Method: "foo"}))
	assert.Equal("[error] only a test\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	log.Fatal(fmt.Errorf("only a test"), OptErrorEventState(&http.Request{Method: "foo"}))
	assert.Equal("[fatal] only a test\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	log.Path = []string{"outer", "inner"}
	log.Labels = Labels{"foo": "bar"}
	log.Info("format test")
	assert.Equal("[outer > inner] [info] format test\tfoo=bar\n", buf.String())
}

func TestScopeFromContext(t *testing.T) {
	assert := assert.New(t)

	sc := NewScope(None())
	sc.Path = []string{"one", "two"}
	sc.Labels = Labels{"foo": "bar"}

	ctx := WithLabels(context.Background(), Labels{"moo": "loo"})
	ctx = WithPath(ctx, "three", "four")

	final := sc.FromContext(ctx)
	assert.Equal([]string{"one", "two", "three", "four"}, final.Path)
	assert.Equal("bar", final.Labels["foo"])
	assert.Equal("loo", final.Labels["moo"])
}

func TestScopeApply(t *testing.T) {
	assert := assert.New(t)

	sc := NewScope(None())
	sc.Path = []string{"one", "two"}
	sc.Labels = Labels{"foo": "bar"}

	ctx := WithLabels(context.Background(), Labels{"moo": "loo"})
	ctx = WithPath(ctx, "three", "four")

	final := sc.ApplyContext(ctx)
	assert.Equal([]string{"one", "two", "three", "four"}, GetPath(final))
	assert.Equal("bar", GetLabels(final)["foo"])
	assert.Equal("loo", GetLabels(final)["moo"])
}

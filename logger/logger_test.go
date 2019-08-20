package logger

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/uuid"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	log, err := New()
	assert.Nil(err)
	assert.NotNil(log.Scope)
	assert.NotNil(log.Formatter)
	assert.NotNil(log.Output)
	assert.True(log.RecoverPanics)

	for _, defaultFlag := range DefaultFlags {
		assert.True(log.Flags.IsEnabled(defaultFlag))
	}

	log, err = New(OptAll(), OptFormatter(NewJSONOutputFormatter()))
	assert.Nil(err)
	assert.True(log.Flags.IsEnabled(uuid.V4().String()))
	typed, ok := log.Formatter.(*JSONOutputFormatter)
	assert.True(ok)
	assert.NotNil(typed)
}

func TestLoggerE2ESubContext(t *testing.T) {
	assert := assert.New(t)

	output := new(bytes.Buffer)
	log, err := New(
		OptOutput(output),
		OptText(OptTextHideTimestamp(), OptTextNoColor()),
	)
	assert.Nil(err)

	scID := uuid.V4().String()
	sc := log.WithPath(scID)

	sc.Infof("this is infof")
	sc.Errorf("this is errorf")
	sc.Fatalf("this is fatalf")

	sc.Trigger(context.Background(), NewMessageEvent(Info, "this is a triggered message"))

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	assert.Nil(log.DrainContext(ctx))

	assert.Contains(output.String(), fmt.Sprintf("[%s] [info] this is infof", scID))
	assert.Contains(output.String(), fmt.Sprintf("[%s] [error] this is errorf", scID))
	assert.Contains(output.String(), fmt.Sprintf("[%s] [fatal] this is fatalf", scID))
	assert.Contains(output.String(), fmt.Sprintf("[%s] [info] this is a triggered message", scID))
}

func TestLoggerE2ESubContextFields(t *testing.T) {
	assert := assert.New(t)

	output := new(bytes.Buffer)
	log, err := New(
		OptOutput(output),
		OptText(OptTextHideTimestamp(), OptTextNoColor()),
	)
	assert.Nil(err)

	fieldKey := uuid.V4().String()
	fieldValue := uuid.V4().String()
	sc := log.WithLabels(Labels{fieldKey: fieldValue})

	sc.Infof("this is infof")
	sc.Errorf("this is errorf")
	sc.Fatalf("this is fatalf")

	sc.Trigger(context.Background(), NewMessageEvent(Info, "this is a triggered message"))
	assert.Nil(log.DrainContext(context.Background()))

	assert.Contains(output.String(), fmt.Sprintf("[info] this is infof\t%s=%s", fieldKey, fieldValue))
	assert.Contains(output.String(), fmt.Sprintf("[error] this is errorf\t%s=%s", fieldKey, fieldValue))
	assert.Contains(output.String(), fmt.Sprintf("[fatal] this is fatalf\t%s=%s", fieldKey, fieldValue))
	assert.Contains(output.String(), fmt.Sprintf("[info] this is a triggered message\t%s=%s", fieldKey, fieldValue))
}

func TestLoggerSkipTrigger(t *testing.T) {
	assert := assert.New(t)

	output := new(bytes.Buffer)
	log, err := New(
		OptOutput(output),
		OptText(OptTextHideTimestamp(), OptTextNoColor()),
	)
	assert.Nil(err)

	var wasCalled bool
	log.Listen(Info, "---", func(ctx context.Context, e Event) {
		wasCalled = true
	})

	log.Trigger(WithSkipTrigger(context.Background()), NewMessageEvent(Info, "this is a triggered message"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	assert.Nil(log.DrainContext(ctx))

	assert.False(wasCalled)
	assert.Contains(output.String(), "[info] this is a triggered message")
}

func TestLoggerSkipWrite(t *testing.T) {
	assert := assert.New(t)

	output := new(bytes.Buffer)
	log, err := New(
		OptOutput(output),
		OptText(OptTextHideTimestamp(), OptTextNoColor()),
	)
	assert.Nil(err)

	var wasCalled bool
	log.Listen(Info, "---", func(ctx context.Context, e Event) {
		wasCalled = true
	})
	assert.True(log.HasListener(Info, "---"))

	log.Trigger(WithSkipWrite(context.Background()), NewMessageEvent(Info, "this is a triggered message"))

	// at the very least this cannot cause a deadlock.
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	assert.Nil(log.DrainContext(ctx))

	assert.True(wasCalled)
	assert.Empty(output.String())
}

func TestLoggerListeners(t *testing.T) {
	assert := assert.New(t)

	log := MustNew()
	defer log.Close()

	assert.Empty(log.Listeners)
	log.Listen(Info, "foo", NewMessageEventListener(func(_ context.Context, me MessageEvent) {}))
	assert.NotEmpty(log.Listeners)
	assert.True(log.HasListeners(Info))
	assert.True(log.HasListener(Info, "foo"))
	assert.False(log.HasListener(Info, "bar"))

	log.Listen(Error, "foo", NewMessageEventListener(func(_ context.Context, me MessageEvent) {}))
	assert.True(log.HasListeners(Error))
	assert.True(log.HasListener(Error, "foo"))
	assert.False(log.HasListener(Error, "bar"))

	log.Listen(Info, "bar", NewMessageEventListener(func(_ context.Context, me MessageEvent) {}))
	assert.True(log.HasListeners(Info))
	assert.True(log.HasListener(Info, "foo"))
	assert.True(log.HasListener(Info, "bar"))

	log.Listen(Error, "bar", NewMessageEventListener(func(_ context.Context, me MessageEvent) {}))
	assert.True(log.HasListeners(Error))
	assert.True(log.HasListener(Error, "foo"))
	assert.True(log.HasListener(Error, "bar"))

	log.RemoveListener(Info, "foo")
	assert.True(log.HasListeners(Info))
	assert.False(log.HasListener(Info, "foo"))
	assert.True(log.HasListener(Info, "bar"))

	log.RemoveListeners(Error)
	assert.False(log.HasListeners(Error))
	assert.False(log.HasListener(Error, "foo"))
	assert.False(log.HasListener(Error, "bar"))
}

func TestLoggerProd(t *testing.T) {
	assert := assert.New(t)

	p := Prod(OptEnabled("bailey"))
	defer p.Close()

	assert.True(p.Flags.IsEnabled("bailey"))
	assert.True(p.Formatter.(*TextOutputFormatter).NoColor)
}

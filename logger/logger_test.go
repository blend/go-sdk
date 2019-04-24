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
	assert.NotNil(log.Context)
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
	sc := log.SubContext(scID)

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
	sc := log.WithFields(Fields{fieldKey: fieldValue})

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

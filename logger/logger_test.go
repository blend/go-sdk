/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

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
	assert.NotNil(log.Flags)
	assert.NotNil(log.Writable)
	assert.NotNil(log.Scope)
	assert.NotNil(log.Formatter)
	assert.NotNil(log.Output)
	assert.True(log.RecoverPanics)

	for _, defaultFlag := range DefaultFlags {
		assert.True(log.Flags.IsEnabled(defaultFlag))
	}
	assert.True(log.Writable.All())

	log, err = New(OptAll(), OptFormatter(NewJSONOutputFormatter()))
	assert.Nil(err)
	assert.True(log.Flags.IsEnabled(uuid.V4().String()))
	typed, ok := log.Formatter.(*JSONOutputFormatter)
	assert.True(ok)
	assert.NotNil(typed)
}

func TestLoggerFlagsWritten(t *testing.T) {
	its := assert.New(t)

	buf := new(bytes.Buffer)
	log := Memory(buf)
	defer log.Close()

	log.Writable.Disable(Info)

	eventTriggered := make(chan struct{})
	log.Listen(Info, DefaultListenerName, func(_ context.Context, e Event) {
		close(eventTriggered)
	})

	log.Dispatch(context.TODO(), NewMessageEvent(Info, "test"))
	<-eventTriggered
	its.Empty(buf.String())

	log.Dispatch(context.TODO(), NewMessageEvent(Error, "this is just a test"))
	its.Equal("[error] this is just a test\n", buf.String())
}

func TestLoggerScopes(t *testing.T) {
	its := assert.New(t)

	buf := new(bytes.Buffer)
	log := Memory(buf)
	defer log.Close()

	log.Scopes.Disable("disabled/*")

	eventTriggered := make(chan struct{})
	log.Listen(Error, DefaultListenerName, func(_ context.Context, e Event) {
		close(eventTriggered)
	})

	log.WithPath("disabled", "foo").TriggerContext(context.TODO(), NewMessageEvent(Info, "test"))
	// should not trigger with scope disabled
	its.Empty(buf.String())

	log.WithPath("not-disabled", "foo").TriggerContext(context.TODO(), NewMessageEvent(Error, "this is just a test"))
	<-eventTriggered
	its.Equal("[not-disabled > foo] [error] this is just a test\n", buf.String())
}

func TestLoggerWritableScopes(t *testing.T) {
	its := assert.New(t)

	buf := new(bytes.Buffer)
	log := Memory(buf)
	defer log.Close()

	log.WritableScopes.Disable("disabled/*")

	eventTriggered := make(chan struct{})
	log.Listen(Error, DefaultListenerName, func(_ context.Context, e Event) {
		close(eventTriggered)
	})

	log.WithPath("disabled", "foo").TriggerContext(context.TODO(), NewMessageEvent(Error, "test"))
	<-eventTriggered
	its.Empty(buf.String())

	eventTriggered = make(chan struct{})
	log.WithPath("not-disabled", "foo").TriggerContext(context.TODO(), NewMessageEvent(Error, "this is just a test"))
	<-eventTriggered
	its.Equal("[not-disabled > foo] [error] this is just a test\n", buf.String())
}

func TestLoggerE2ESubContext(t *testing.T) {
	assert := assert.New(t)

	output := new(bytes.Buffer)
	log, err := New(
		OptOutput(output),
		OptText(OptTextHideTimestamp(), OptTextNoColor()),
	)
	assert.Nil(err)
	defer log.Close()

	scID := uuid.V4().String()
	sc := log.WithPath(scID)

	sc.Infof("this is infof")
	sc.Errorf("this is errorf")
	sc.Fatalf("this is fatalf")

	sc.Trigger(NewMessageEvent(Info, "this is a triggered message"))

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	log.DrainContext(ctx)

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

	sc.Trigger(NewMessageEvent(Info, "this is a triggered message"))
	log.DrainContext(context.Background())

	assert.Contains(output.String(), fmt.Sprintf("[info] this is infof\t%s=%s", fieldKey, fieldValue))
	assert.Contains(output.String(), fmt.Sprintf("[error] this is errorf\t%s=%s", fieldKey, fieldValue))
	assert.Contains(output.String(), fmt.Sprintf("[fatal] this is fatalf\t%s=%s", fieldKey, fieldValue))
	assert.Contains(output.String(), fmt.Sprintf("[info] this is a triggered message\t%s=%s", fieldKey, fieldValue))
}

func TestLoggerAnnotations(t *testing.T) {
	assert := assert.New(t)

	output := new(bytes.Buffer)
	log, err := New(
		OptOutput(output),
		OptText(OptTextHideTimestamp(), OptTextNoColor()),
	)
	assert.Nil(err)

	fieldKey := uuid.V4().String()
	fieldValue := uuid.V4().String()
	sc := log.WithAnnotations(Annotations{fieldKey: fieldValue})

	sc.Infof("this is infof")
	sc.Errorf("this is errorf")
	sc.Fatalf("this is fatalf")

	sc.Trigger(NewMessageEvent(Info, "this is a triggered message"))
	log.DrainContext(context.Background())

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
	defer log.Close()

	var wasCalled bool
	log.Listen(Info, "---", func(ctx context.Context, e Event) {
		wasCalled = true
	})

	log.TriggerContext(WithSkipTrigger(context.Background(), true), NewMessageEvent(Info, "this is a triggered message"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	log.DrainContext(ctx)

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
	defer log.Close()

	var wasCalled bool
	log.Listen(Info, "---", func(ctx context.Context, e Event) {
		wasCalled = true
	})
	assert.True(log.HasListener(Info, "---"))

	log.TriggerContext(WithSkipWrite(context.Background(), true), NewMessageEvent(Info, "this is a triggered message"))

	// at the very least this cannot cause a deadlock.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.DrainContext(ctx)
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

	assert.Nil(log.RemoveListener(Info, "foo"))
	assert.True(log.HasListeners(Info))
	assert.False(log.HasListener(Info, "foo"))
	assert.True(log.HasListener(Info, "bar"))

	assert.Nil(log.RemoveListeners(Error))
	assert.False(log.HasListeners(Error))
	assert.False(log.HasListener(Error, "foo"))
	assert.False(log.HasListener(Error, "bar"))
}

func TestLoggerFilters(t *testing.T) {
	assert := assert.New(t)

	log := MustNew()
	defer log.Close()

	noop := func(_ context.Context, e MessageEvent) (MessageEvent, bool) {
		return e, false
	}

	assert.Empty(log.Filters)
	log.Filter(Info, "foo", NewMessageEventFilter(noop))
	assert.NotEmpty(log.Filters)
	assert.True(log.HasFilters(Info))
	assert.True(log.HasFilter(Info, "foo"))
	assert.False(log.HasFilter(Info, "bar"))

	log.Filter(Error, "foo", NewMessageEventFilter(noop))
	assert.True(log.HasFilters(Error))
	assert.True(log.HasFilter(Error, "foo"))
	assert.False(log.HasFilter(Error, "bar"))

	log.Filter(Info, "bar", NewMessageEventFilter(noop))
	assert.True(log.HasFilters(Info))
	assert.True(log.HasFilter(Info, "foo"))
	assert.True(log.HasFilter(Info, "bar"))

	log.Filter(Error, "bar", NewMessageEventFilter(noop))
	assert.True(log.HasFilters(Error))
	assert.True(log.HasFilter(Error, "foo"))
	assert.True(log.HasFilter(Error, "bar"))

	log.RemoveFilter(Info, "foo")
	assert.True(log.HasFilters(Info))
	assert.False(log.HasFilter(Info, "foo"))
	assert.True(log.HasFilter(Info, "bar"))

	log.RemoveFilters(Error)
	assert.False(log.HasFilters(Error))
	assert.False(log.HasFilter(Error, "foo"))
	assert.False(log.HasFilter(Error, "bar"))
}

func TestLoggerDispatchFilterMutate(t *testing.T) {
	it := assert.New(t)

	output := new(bytes.Buffer)
	log, err := New(
		OptOutput(output),
		OptText(OptTextHideTimestamp(), OptTextNoColor()),
	)
	it.Nil(err)
	defer log.Close()

	var wasCalled bool
	var textWasModified bool

	log.Listen(Info, "---", func(ctx context.Context, e Event) {
		wasCalled = true
		textWasModified = e.(MessageEvent).Text == "not_test_message"
	})
	var wasFiltered bool
	log.Filter(Info, "---", func(ctx context.Context, e Event) (Event, bool) {
		wasFiltered = true
		copy := e.(MessageEvent)
		copy.Text = "not_" + copy.Text
		return copy, false
	})

	log.TriggerContext(context.Background(), NewMessageEvent(Info, "test_message"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	log.DrainContext(ctx)

	it.True(wasFiltered)
	it.True(wasCalled)
	it.True(textWasModified)

	it.Equal("[info] not_test_message\n", output.String())
}

func TestLoggerDispatchFilterDrop(t *testing.T) {
	it := assert.New(t)

	output := new(bytes.Buffer)
	log, err := New(
		OptOutput(output),
		OptText(OptTextHideTimestamp(), OptTextNoColor()),
	)
	it.Nil(err)
	defer log.Close()

	var wasCalled bool
	log.Listen(Info, "---", func(ctx context.Context, e Event) {
		wasCalled = true
	})
	var wasFiltered bool
	log.Filter(Info, "---", func(ctx context.Context, e Event) (Event, bool) {
		wasFiltered = true
		return nil, true
	})

	log.TriggerContext(context.Background(), NewMessageEvent(Info, "this is a triggered message"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	log.DrainContext(ctx)

	it.True(wasFiltered)
	it.False(wasCalled)
	it.Empty(output.String())
}

func TestLoggerDrain(t *testing.T) {
	assert := assert.New(t)

	log := MustNew()
	defer log.Close()
	assert.Empty(log.Listeners)

	eventsCounted := 0
	log.Listen(Info, "foo", NewMessageEventListener(func(_ context.Context, me MessageEvent) {
		eventsCounted++
	}))

	for i := 0; i < 5; i++ {
		log.Info("event")
	}
	log.Drain()
	assert.Equal(5, eventsCounted)

	for i := 0; i < 4; i++ {
		log.Info("event")
	}
	log.Drain()
	assert.Equal(9, eventsCounted)
}

func TestLoggerProd(t *testing.T) {
	assert := assert.New(t)

	p := Prod(OptEnabled("example-string"))
	defer p.Close()

	assert.True(p.Flags.IsEnabled("example-string"))
	assert.True(p.Formatter.(*TextOutputFormatter).NoColor)
}

func TestLoggerClose(t *testing.T) {
	assert := assert.New(t)

	buf0 := new(bytes.Buffer)
	l0 := Memory(buf0)
	buf1 := new(bytes.Buffer)
	l1 := Memory(buf1)

	l0.Close()
	l1.Infof("this is a test")
	assert.Empty(buf0.String())
	assert.NotEmpty(buf1.String())
}

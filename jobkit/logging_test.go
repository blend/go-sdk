package jobkit

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/logger"
)

func TestLoggingDebugf(t *testing.T) {
	assert := assert.New(t)

	// should not panic
	Debugf(nil, nil, "foo")

	ctx := cron.WithJobInvocation(context.Background(), &cron.JobInvocation{
		ID:      "log-test-0",
		JobName: "log-test",
	})

	triggered := make(chan struct{})
	log := logger.All(logger.OptOutput(ioutil.Discard))
	defer log.Close()

	var message string
	log.Listen(logger.Debug, "check-listener", logger.NewMessageEventListener(func(_ context.Context, me logger.MessageEvent) {
		defer close(triggered)
		message = me.Text
	}))

	Debugf(ctx, log, "foo %s", "bar")
	<-triggered
	assert.Equal("foo bar", message)
}

func TestLoggingInfof(t *testing.T) {
	assert := assert.New(t)

	// should not panic
	Infof(nil, nil, "foo")

	ctx := cron.WithJobInvocation(context.Background(), &cron.JobInvocation{
		ID:      "log-test-0",
		JobName: "log-test",
	})

	triggered := make(chan struct{})
	log := logger.All(logger.OptOutput(ioutil.Discard))
	defer log.Close()

	var message string
	log.Listen(logger.Info, "check-listener", logger.NewMessageEventListener(func(_ context.Context, me logger.MessageEvent) {
		defer func() { close(triggered) }()
		message = me.Text
	}))

	Infof(ctx, log, "foo %s", "bar")

	<-triggered
	assert.Equal("foo bar", message)
}

func TestLoggingWarningf(t *testing.T) {
	assert := assert.New(t)

	// should not panic
	Warningf(nil, nil, "foo")

	ctx := cron.WithJobInvocation(context.Background(), &cron.JobInvocation{
		ID:      "log-test-0",
		JobName: "log-test",
	})

	triggered := make(chan struct{})
	log := logger.All(logger.OptOutput(ioutil.Discard))
	defer log.Close()

	var message string
	log.Listen(logger.Warning, "check-listener", logger.NewErrorEventListener(func(_ context.Context, ee logger.ErrorEvent) {
		defer func() { close(triggered) }()
		message = ee.Err.Error()
	}))

	Warningf(ctx, log, "foo %s", "bar")

	<-triggered
	assert.Equal("foo bar", message)
}

func TestLoggingWarning(t *testing.T) {
	assert := assert.New(t)

	// should not panic
	Warning(nil, nil, fmt.Errorf("foo"))

	ctx := cron.WithJobInvocation(context.Background(), &cron.JobInvocation{
		ID:      "log-test-0",
		JobName: "log-test",
	})

	triggered := make(chan struct{})
	log := logger.All(logger.OptOutput(ioutil.Discard))
	defer log.Close()

	var message string
	log.Listen(logger.Warning, "check-listener", logger.NewErrorEventListener(func(_ context.Context, ee logger.ErrorEvent) {
		defer func() { close(triggered) }()
		message = ee.Err.Error()
	}))

	Warning(ctx, log, fmt.Errorf("foo %s", "bar"))

	<-triggered
	assert.Equal("foo bar", message)
}

func TestLoggingErrorf(t *testing.T) {
	assert := assert.New(t)

	// should not panic
	Errorf(nil, nil, "foo")

	ctx := cron.WithJobInvocation(context.Background(), &cron.JobInvocation{
		ID:      "log-test-0",
		JobName: "log-test",
	})

	triggered := make(chan struct{})
	log := logger.All(logger.OptOutput(ioutil.Discard))
	defer log.Close()

	var message string
	log.Listen(logger.Error, "check-listener", logger.NewErrorEventListener(func(_ context.Context, ee logger.ErrorEvent) {
		defer func() { close(triggered) }()
		message = ee.Err.Error()
	}))

	Errorf(ctx, log, "foo %s", "bar")

	<-triggered
	assert.Equal("foo bar", message)
}

func TestLoggingError(t *testing.T) {
	assert := assert.New(t)

	// should not panic
	Error(nil, nil, fmt.Errorf("foo"))

	ctx := cron.WithJobInvocation(context.Background(), &cron.JobInvocation{
		ID:      "log-test-0",
		JobName: "log-test",
	})

	triggered := make(chan struct{})
	log := logger.All(logger.OptOutput(ioutil.Discard))
	defer log.Close()

	var message string
	log.Listen(logger.Error, "check-listener", logger.NewErrorEventListener(func(_ context.Context, ee logger.ErrorEvent) {
		defer func() { close(triggered) }()
		message = ee.Err.Error()
	}))

	Error(ctx, log, fmt.Errorf("foo %s", "bar"))

	<-triggered
	assert.Equal("foo bar", message)
}

func TestLoggingFatalf(t *testing.T) {
	assert := assert.New(t)

	// should not panic
	Fatalf(nil, nil, "foo")

	ctx := cron.WithJobInvocation(context.Background(), &cron.JobInvocation{
		ID:      "log-test-0",
		JobName: "log-test",
	})

	triggered := make(chan struct{})
	log := logger.All(logger.OptOutput(ioutil.Discard))
	defer log.Close()

	var message string
	log.Listen(logger.Fatal, "check-listener", logger.NewErrorEventListener(func(_ context.Context, ee logger.ErrorEvent) {
		defer func() { close(triggered) }()
		message = ee.Err.Error()
	}))

	Fatalf(ctx, log, "foo %s", "bar")

	<-triggered
	assert.Equal("foo bar", message)
}

func TestLoggingFatal(t *testing.T) {
	assert := assert.New(t)

	// should not panic
	Fatal(nil, nil, fmt.Errorf("foo"))

	ctx := cron.WithJobInvocation(context.Background(), &cron.JobInvocation{
		ID:      "log-test-0",
		JobName: "log-test",
	})

	triggered := make(chan struct{})
	log := logger.All(logger.OptOutput(ioutil.Discard))
	defer log.Close()

	var message string
	log.Listen(logger.Fatal, "check-listener", logger.NewErrorEventListener(func(_ context.Context, ee logger.ErrorEvent) {
		defer func() { close(triggered) }()
		message = ee.Err.Error()
	}))

	Fatal(ctx, log, fmt.Errorf("foo %s", "bar"))

	<-triggered
	assert.Equal("foo bar", message)
}

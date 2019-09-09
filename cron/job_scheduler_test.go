package cron

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/uuid"
)

var (
	_ graceful.Graceful = (*JobScheduler)(nil)
)

func TestJobSchedulerCullHistoryMaxAge(t *testing.T) {
	assert := assert.New(t)

	js := NewJobScheduler(NewJob("foo", noop))
	js.HistoryMaxCountProvider = func() int { return 10 }
	js.HistoryMaxAgeProvider = func() time.Duration { return 6 * time.Hour }

	js.History = []JobInvocation{
		{ID: uuid.V4().String(), Started: time.Now().Add(-10 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-9 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-8 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-7 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-6 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-5 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-4 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-3 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-2 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-1 * time.Hour)},
	}

	filtered := js.cullHistory()
	assert.Len(filtered, 5)
}

func TestJobSchedulerCullHistoryMaxCount(t *testing.T) {
	assert := assert.New(t)

	js := NewJobScheduler(NewJob("foo", noop))
	js.HistoryMaxCountProvider = func() int { return 5 }
	js.HistoryMaxAgeProvider = func() time.Duration { return 6 * time.Hour }

	js.History = []JobInvocation{
		{ID: uuid.V4().String(), Started: time.Now().Add(-10 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-9 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-8 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-7 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-6 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-5 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-4 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-3 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-2 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-1 * time.Minute)},
	}

	filtered := js.cullHistory()
	assert.Len(filtered, 5)
}

func TestJobSchedulerEnableDisable(t *testing.T) {
	assert := assert.New(t)

	var enabled, disabled bool

	js := NewJobScheduler(
		NewJob("foo",
			noop,
			OptJobBuilderOnDisabled(func(_ context.Context) { disabled = true }),
			OptJobBuilderOnEnabled(func(_ context.Context) { enabled = true }),
		),
	)

	js.HistoryMaxCountProvider = func() int { return 5 }
	js.HistoryMaxAgeProvider = func() time.Duration { return 6 * time.Hour }

	js.Disable()
	assert.True(js.Disabled)

	assert.False(js.CanScheduledRun())

	js.Enable()
	assert.False(js.Disabled)

	assert.True(disabled)
	assert.True(enabled)
}

func TestJobSchedulerPersistHistory(t *testing.T) {
	assert := assert.New(t)

	js := NewJobScheduler(
		NewJob("foo", noop),
	)
	js.HistoryEnabledProvider = func() bool { return true }

	assert.Nil(js.RestoreHistory(context.Background()))
	assert.Nil(js.PersistHistory(context.Background()))

	var history []JobInvocation
	js.HistoryPersistProvider = func(_ context.Context, h []JobInvocation) error {
		history = h
		return nil
	}
	js.Run()
	assert.Len(history, 1)
	js.Run()
	assert.Len(history, 2)

	js.HistoryEnabledProvider = func() bool { return false }
	js.Run()
	assert.Len(history, 2)

	js.HistoryRestoreProvider = func(_ context.Context) ([]JobInvocation, error) {
		return []JobInvocation{
			*NewJobInvocation("foo"),
			*NewJobInvocation("foo"),
			*NewJobInvocation("foo"),
		}, nil
	}
	assert.Nil(js.RestoreHistory(context.Background()))
	assert.Len(js.History, 3)

	js.HistoryPersistProvider = func(_ context.Context, h []JobInvocation) error {
		return fmt.Errorf("only a test")
	}
	assert.NotNil(js.PersistHistory(context.Background()))
}

func TestJobSchedulerLabels(t *testing.T) {
	assert := assert.New(t)

	js := NewJobScheduler(NewJob("test", noop))
	labels := js.Labels()
	assert.Equal("test", labels["name"])

	js.LabelsProvider = func() map[string]string {
		return map[string]string{
			"name": "not-test",
			"foo":  "bar",
			"fuzz": "wuzz",
		}
	}

	labels = js.Labels()
	assert.Equal("not-test", labels["name"])
	assert.Equal("bar", labels["foo"])
	assert.Equal("wuzz", labels["fuzz"])
}

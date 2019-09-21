package jobkit

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/ref"
	"github.com/blend/go-sdk/slack"
	"github.com/blend/go-sdk/uuid"
)

func scheduleProvider(schedule cron.Schedule) func() cron.Schedule {
	return func() cron.Schedule {
		return schedule
	}
}

func TestNewJob(t *testing.T) {
	assert := assert.New(t)

	cfg := JobConfig{
		JobConfig: cron.JobConfig{
			Name:    "test",
			Timeout: 2 * time.Second,
		},
		Schedule: "@every 1s",
	}
	var didCallAction bool
	action := func(_ context.Context) error {
		didCallAction = true
		return nil
	}

	job, err := NewJob(cfg, action, OptParsedSchedule("@every 2s"))
	assert.Nil(err)
	assert.NotNil(job)
	assert.Equal(2*time.Second, job.CompiledSchedule.(cron.IntervalSchedule).Every)
	assert.Nil(job.Action(context.Background()))
	assert.True(didCallAction)
}

func TestJobWrap(t *testing.T) {
	assert := assert.New(t)

	job := cron.NewJob(cron.OptJobName("test"), cron.OptJobAction(func(_ context.Context) error {
		return nil
	}))
	wrapped := Wrap(job)
	assert.Equal(wrapped.Name(), job.Name())
	assert.NotNil(job.Action)
}

func TestJobLifecycleHooksNotificationsUnset(t *testing.T) {
	assert := assert.New(t)

	ctx := cron.WithJobInvocation(context.Background(), &cron.JobInvocation{
		ID:      uuid.V4().String(),
		JobName: "test-job",
	})

	slackMessages := make(chan slack.Message, 16)

	job := &Job{
		Config:      JobConfig{},
		SlackClient: slack.MockWebhookSender(slackMessages),
	}

	job.OnStart(ctx)
	assert.Empty(slackMessages)

	job.OnComplete(ctx)
	assert.Empty(slackMessages)

	job.OnFailure(ctx)
	assert.NotEmpty(slackMessages)

	job.OnCancellation(ctx)
	assert.NotEmpty(slackMessages)

	job.OnBroken(ctx)
	assert.NotEmpty(slackMessages)

	job.OnFixed(ctx)
	assert.NotEmpty(slackMessages)
}

func TestJobLifecycleHooksNotificationsSetDisabled(t *testing.T) {
	assert := assert.New(t)

	ctx := cron.WithJobInvocation(context.Background(), &cron.JobInvocation{
		ID:      uuid.V4().String(),
		JobName: "test-job",
	})

	slackMessages := make(chan slack.Message, 1)

	job := &Job{
		SlackClient: slack.MockWebhookSender(slackMessages),
		Config: JobConfig{
			NotifyOnStart:        ref.Bool(false),
			NotifyOnSuccess:      ref.Bool(false),
			NotifyOnFailure:      ref.Bool(false),
			NotifyOnBroken:       ref.Bool(false),
			NotifyOnFixed:        ref.Bool(false),
			NotifyOnCancellation: ref.Bool(false),
		},
	}

	job.OnStart(ctx)
	assert.Empty(slackMessages)

	job.OnComplete(ctx)
	assert.Empty(slackMessages)

	job.OnFailure(ctx)
	assert.Empty(slackMessages)

	job.OnCancellation(ctx)
	assert.Empty(slackMessages)

	job.OnBroken(ctx)
	assert.Empty(slackMessages)

	job.OnFixed(ctx)
	assert.Empty(slackMessages)
}

func TestJobLifecycleHooksNotificationsSetEnabled(t *testing.T) {
	assert := assert.New(t)

	ctx := cron.WithJobInvocation(context.Background(), &cron.JobInvocation{
		ID:      uuid.V4().String(),
		JobName: "test-job",
		Err:     fmt.Errorf("only a test"),
	})

	slackMessages := make(chan slack.Message, 6)

	job := &Job{
		SlackClient: slack.MockWebhookSender(slackMessages),
		Config: JobConfig{
			NotifyOnStart:        ref.Bool(true),
			NotifyOnSuccess:      ref.Bool(true),
			NotifyOnFailure:      ref.Bool(true),
			NotifyOnBroken:       ref.Bool(true),
			NotifyOnFixed:        ref.Bool(true),
			NotifyOnCancellation: ref.Bool(true),
		},
	}

	job.OnStart(ctx)
	job.OnComplete(ctx)
	job.OnFailure(ctx)
	job.OnCancellation(ctx)
	job.OnBroken(ctx)
	job.OnFixed(ctx)

	assert.Len(slackMessages, 6)

	msg := <-slackMessages
	assert.Contains("cron.started", msg.Text)

	msg = <-slackMessages
	assert.Contains("cron.complete", msg.Text)

	msg = <-slackMessages
	assert.Contains("cron.failed", msg.Text)

	msg = <-slackMessages
	assert.Contains("cron.cancelled", msg.Text)

	msg = <-slackMessages
	assert.Contains("cron.broken", msg.Text)

	msg = <-slackMessages
	assert.Contains("cron.fixed", msg.Text)
}

func TestJobHistoryProvider(t *testing.T) {
	assert := assert.New(t)

	tmpdir, err := ioutil.TempDir("", "gosdk_jobkit")
	assert.Nil(err)
	defer os.RemoveAll(tmpdir)

	job := &Job{
		Config: JobConfig{
			HistoryPath: tmpdir,
		},
	}

	history := []cron.JobInvocation{
		createTestCompleteJobInvocation("test0", 100*time.Millisecond),
		createTestCompleteJobInvocation("test0", 200*time.Millisecond),
		createTestFailedJobInvocation("test0", 100*time.Millisecond, fmt.Errorf("this is only a test")),
	}

	err = job.PersistHistory(context.Background(), history)
	assert.Nil(err)

	returned, err := job.RestoreHistory(context.Background())
	assert.Nil(err)
	assert.Len(returned, 3)
}

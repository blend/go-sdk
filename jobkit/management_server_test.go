package jobkit

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/web"
)

func TestManagementServer(t *testing.T) {
	assert := assert.New(t)

	jm := cron.New()

	jm.LoadJob(cron.NewJob("test0", func(_ context.Context) error { return nil }))
	jm.LoadJob(cron.NewJob("test1", func(_ context.Context) error { return nil }))

	app := NewManagementServer(jm, &Config{
		Web: web.Config{
			Port: 5000,
		},
	})

	meta, err := app.Mock().Get("/").ExecuteWithMeta()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)

	var jobs cron.Status
	meta, err = app.Mock().Get("/api/jobs").JSONWithMeta(&jobs)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Len(jobs.Jobs, 2)
}

func TestManagementServerHealthz(t *testing.T) {
	assert := assert.New(t)

	jm := cron.New()
	jm.LoadJob(cron.NewJob("test0", func(_ context.Context) error { return nil }))
	jm.LoadJob(cron.NewJob("test1", func(_ context.Context) error { return nil }))
	jm.Start()

	app := NewManagementServer(jm, &Config{
		Web: web.Config{
			Port: 5000,
		},
	})

	meta, err := app.Mock().Get("/healthz").ExecuteWithMeta()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)

	jm.Stop()

	meta, err = app.Mock().Get("/healthz").ExecuteWithMeta()
	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, meta.StatusCode)
}

func TestManagementServerIndex(t *testing.T) {
	assert := assert.New(t)

	jobName := "test0"
	invocationID := uuid.V4().String()
	output := uuid.V4().String()
	errorOutput := uuid.V4().String()

	jm := cron.New()
	jm.LoadJob(cron.NewJob(jobName, func(_ context.Context) error { return nil }))

	js, err := jm.Job(jobName)
	assert.Nil(err)
	js.History = []cron.JobInvocation{
		{
			ID:      invocationID,
			JobName: jobName,
			State: JobInvocationState{
				Output:      bytes.NewBufferString(output),
				ErrorOutput: bytes.NewBufferString(errorOutput),
			},
		},
	}

	app := NewManagementServer(jm, &Config{
		Web: web.Config{
			Port: 5000,
		},
	})

	contents, meta, err := app.Mock().Get("/").BytesWithMeta()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Contains(string(contents), jobName)
	assert.Contains(string(contents), invocationID)

	contents, meta, err = app.Mock().Get(fmt.Sprintf("/job.invocation/%s/%s", jobName, invocationID)).BytesWithMeta()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Contains(string(contents), jobName)
	assert.Contains(string(contents), invocationID)
	assert.Contains(string(contents), output)
	assert.Contains(string(contents), errorOutput)
}

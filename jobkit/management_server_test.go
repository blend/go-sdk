package jobkit

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/stringutil"
	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/web"
)

func TestManagementServer(t *testing.T) {
	assert := assert.New(t)

	jm := cron.New()

	jm.LoadJobs(
		cron.NewJob("test0", func(_ context.Context) error { return nil }),
		cron.NewJob("test1", func(_ context.Context) error { return nil }),
	)

	app := NewManagementServer(jm, Config{
		Web: web.Config{
			Port: 5000,
		},
	})

	meta, err := web.MockGet(app, "/").Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)

	var jobs cron.Status
	meta, err = web.MockGet(app, "/api/jobs").JSON(&jobs)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Len(jobs.Jobs, 2)
}

func TestManagementServerHealthz(t *testing.T) {
	assert := assert.New(t)

	jm := cron.New()
	jm.LoadJobs(
		cron.NewJob("test0", func(_ context.Context) error { return nil }),
		cron.NewJob("test1", func(_ context.Context) error { return nil }),
	)
	jm.StartAsync()
	app := NewManagementServer(jm, Config{
		Web: web.Config{
			Port: 5000,
		},
	})

	meta, err := web.MockGet(app, "/healthz").Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)

	jm.Stop()

	meta, err = web.MockGet(app, "/healthz").Discard()
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
	jm.LoadJobs(cron.NewJob(jobName, func(_ context.Context) error { return nil }))

	js, err := jm.Job(jobName)
	assert.Nil(err)
	js.History = []cron.JobInvocation{
		{
			ID:      invocationID,
			JobName: jobName,
			Output:  stringutil.NewLineBuffer(bytes.NewBufferString(output + errorOutput).Bytes()),
		},
	}

	app := NewManagementServer(jm, Config{
		Web: web.Config{
			Port: 5000,
		},
	})

	contents, meta, err := web.MockGet(app, "/").Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Contains(string(contents), jobName)
	assert.Contains(string(contents), invocationID)

	contents, meta, err = web.MockGet(app, fmt.Sprintf("/job.invocation/%s/%s", jobName, invocationID)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Contains(string(contents), jobName)
	assert.Contains(string(contents), invocationID)
	assert.Contains(string(contents), output)
	assert.Contains(string(contents), errorOutput)
}

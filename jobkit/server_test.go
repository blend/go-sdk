package jobkit

import (
	"bytes"
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

	jm.LoadJobs(
		cron.NewJob(cron.OptJobName("test0")),
		cron.NewJob(cron.OptJobName("test1")),
	)

	app := NewServer(jm, Config{
		Web: web.Config{
			Port: 5000,
		},
	})

	meta, err := web.MockGet(app, "/").Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)

	var jobs []*cron.JobScheduler
	meta, err = web.MockGet(app, "/api/jobs").JSON(&jobs)

	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Len(jobs, 2)
}

func TestManagementServerStatus(t *testing.T) {
	assert := assert.New(t)

	jm := cron.New()
	jm.LoadJobs(
		cron.NewJob(cron.OptJobName("test0")),
		cron.NewJob(cron.OptJobName("test1")),
	)
	jm.StartAsync()
	app := NewServer(jm, Config{
		Web: web.Config{
			Port: 5000,
		},
	})

	var status cron.JobManagerStatus
	meta, err := web.MockGet(app, "/status.json").JSON(&status)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Equal(cron.JobManagerStateRunning, status.State)

	jm.Stop()

	meta, err = web.MockGet(app, "/status.json").JSON(&status)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Equal(cron.JobManagerStateStopped, status.State)
}

func TestManagementServerIndex(t *testing.T) {
	assert := assert.New(t)

	jobName := "test0"
	invocationID := uuid.V4().String()
	output := uuid.V4().String()
	errorOutput := uuid.V4().String()

	jm := cron.New()
	jm.LoadJobs(cron.NewJob(cron.OptJobName(jobName)))

	js, err := jm.Job(jobName)
	assert.Nil(err)
	js.History = []cron.JobInvocation{
		{
			ID:      invocationID,
			JobName: jobName,
			Output:  cron.NewOutputBuffer(bytes.NewBufferString(output + errorOutput).Bytes()),
		},
	}

	app := NewServer(jm, Config{
		Web: web.Config{
			Port: 5000,
		},
	})

	contents, meta, err := web.MockGet(app, "/").Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Contains(string(contents), jobName)
}

func TestManagementServerJob(t *testing.T) {
	assert := assert.New(t)

	jobName := "test0"
	invocationID := uuid.V4().String()
	output := uuid.V4().String()
	errorOutput := uuid.V4().String()

	jm := cron.New()
	jm.LoadJobs(cron.NewJob(cron.OptJobName(jobName)))

	js, err := jm.Job(jobName)
	assert.Nil(err)
	js.History = []cron.JobInvocation{
		{
			ID:      invocationID,
			JobName: jobName,
			Output:  cron.NewOutputBuffer(bytes.NewBufferString(output + errorOutput).Bytes()),
		},
	}

	app := NewServer(jm, Config{
		Web: web.Config{
			Port: 5000,
		},
	})

	contents, meta, err := web.MockGet(app, fmt.Sprintf("/job/%s", jobName)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Contains(string(contents), jobName)
	assert.Contains(string(contents), invocationID)
}

func TestManagementServerJobInvocation(t *testing.T) {
	assert := assert.New(t)

	jobName := "test0"
	invocationID := uuid.V4().String()
	output := uuid.V4().String()
	errorOutput := uuid.V4().String()

	jm := cron.New()
	jm.LoadJobs(cron.NewJob(cron.OptJobName(jobName)))

	js, err := jm.Job(jobName)
	assert.Nil(err)
	js.History = []cron.JobInvocation{
		{
			ID:      invocationID,
			JobName: jobName,
			Output:  cron.NewOutputBuffer(bytes.NewBufferString(output + errorOutput).Bytes()),
		},
	}

	app := NewServer(jm, Config{
		Web: web.Config{
			Port: 5000,
		},
	})

	contents, meta, err := web.MockGet(app, fmt.Sprintf("/job.invocation/%s/%s", jobName, invocationID)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Contains(string(contents), jobName)
	assert.Contains(string(contents), invocationID)
}

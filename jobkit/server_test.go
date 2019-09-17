package jobkit

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/web"
)

func createTestManagementServer() (*cron.JobManager, *web.App) {

	test0 := cron.NewJob(cron.OptJobName("test0"))
	test1 := cron.NewJob(cron.OptJobName("test1"))

	jm := cron.New()
	jm.LoadJobs(test0, test1)

	jm.Jobs["test0"].History = []cron.JobInvocation{
		{
			ID:       uuid.V4().String(),
			JobName:  "test0",
			Started:  time.Now().UTC(),
			Finished: time.Now().UTC().Add(time.Second),
			State:    cron.JobInvocationStateComplete,
			Elapsed:  time.Second,
		},
		{
			ID:       uuid.V4().String(),
			JobName:  "test0",
			Started:  time.Now().UTC(),
			Finished: time.Now().UTC().Add(time.Second),
			State:    cron.JobInvocationStateComplete,
			Elapsed:  time.Second,
		},
	}
	jm.Jobs["test1"].History = []cron.JobInvocation{
		{
			ID:       uuid.V4().String(),
			JobName:  "test1",
			Started:  time.Now().UTC(),
			Finished: time.Now().UTC().Add(time.Second),
			State:    cron.JobInvocationStateComplete,
			Elapsed:  time.Second,
		},
		{
			ID:       uuid.V4().String(),
			JobName:  "test1",
			Started:  time.Now().UTC(),
			Finished: time.Now().UTC().Add(time.Second),
			State:    cron.JobInvocationStateComplete,
			Elapsed:  time.Second,
		},
	}

	return jm, NewServer(jm, Config{
		Web: web.Config{
			Port: 5000,
		},
	}, web.OptLog(logger.All()))
}

func jobs(jm *cron.JobManager) []*cron.JobScheduler {
	var output []*cron.JobScheduler
	for _, js := range jm.Jobs {
		output = append(output, js)
	}
	return output
}

func TestManagementServerStatus(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()
	jm.StartAsync()

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

	_, app := createTestManagementServer()

	meta, err := web.MockGet(app, "/").Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
}

func TestManagementServerJob(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()

	job := jobs(jm)[0]
	jobName := job.Name()
	invocationID := job.History[0].ID

	contents, meta, err := web.MockGet(app, fmt.Sprintf("/job/%s", jobName)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Contains(string(contents), jobName)
	assert.Contains(string(contents), invocationID)
}

func TestManagementServerJobInvocation(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()

	job := jobs(jm)[0]
	jobName := job.Name()
	invocationID := job.History[0].ID

	contents, meta, err := web.MockGet(app, fmt.Sprintf("/job.invocation/%s/%s", jobName, invocationID)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Contains(string(contents), jobName)
	assert.Contains(string(contents), invocationID)
}

func TestManagementServerPause(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()

	meta, err := web.MockGet(app, "/pause").Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)

	assert.True(jm.Latch.IsPaused())
}

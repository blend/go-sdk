package jobkit

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/web"
)

func TestManagementServerStatus(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()
	jm.StartAsync()
	defer jm.Stop()

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

	jm, app := createTestManagementServer()
	contents, meta, err := web.MockGet(app, "/").Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	jobName := firstJob(jm).Name()
	assert.Contains(string(contents), fmt.Sprintf("/job/%s", jobName))
	assert.Contains(string(contents), "Show job stats and history")
}

func TestManagementServerAPIJobs(t *testing.T) {
	assert := assert.New(t)

	_, app := createTestManagementServer()
	var jobs []cron.JobSchedulerStatus
	meta, err := web.MockGet(app, "/api/jobs").JSON(&jobs)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.NotEmpty(jobs)
}

func TestManagementServerJob(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()

	job := firstJob(jm)
	assert.NotNil(job)
	jobName := job.Name()
	invocationID := job.History[0].ID

	contents, meta, err := web.MockGet(app, fmt.Sprintf("/job/%s", jobName)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Contains(string(contents), jobName)
	assert.Contains(string(contents), invocationID)

	assert.Contains(string(contents), fmt.Sprintf("/job/%s", jobName))
	assert.NotContains(string(contents), "Show job stats and history")
}

func TestManagementServerAPIJob(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()

	job := firstJob(jm)
	assert.NotNil(job)
	jobName := job.Name()

	var js cron.JobSchedulerStatus
	meta, err := web.MockGet(app, fmt.Sprintf("/api/job/%s", jobName)).JSON(&js)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Equal(jobName, js.Name)
}

func TestManagementServerJobNotFound(t *testing.T) {
	assert := assert.New(t)

	_, app := createTestManagementServer()

	meta, err := web.MockGet(app, fmt.Sprintf("/job/%s", uuid.V4().String())).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusNotFound, meta.StatusCode)
}

func TestManagementServerJobInvocation(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()

	job := firstJob(jm)
	assert.NotNil(job)

	jobName := job.Name()
	invocationID := job.History[0].ID

	contents, meta, err := web.MockGet(app, fmt.Sprintf("/job.invocation/%s/%s", jobName, invocationID)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode, string(contents))
	assert.Contains(string(contents), jobName)
	assert.Contains(string(contents), invocationID)
}

func TestManagementServerJobInvocationCurrent(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()

	job := firstJob(jm)
	assert.NotNil(job)

	jobName := job.Name()
	invocationID := job.Current.ID

	contents, meta, err := web.MockGet(app, fmt.Sprintf("/job.invocation/%s/current", jobName)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode, string(contents))
	assert.Contains(string(contents), jobName)
	assert.Contains(string(contents), invocationID)
}

func TestManagementServerJobInvocationNotFound(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()

	meta, err := web.MockGet(app, fmt.Sprintf("/job.invocation/%s/%s", uuid.V4().String(), uuid.V4().String())).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusNotFound, meta.StatusCode)

	job := firstJob(jm)
	assert.NotNil(job)
	jobName := job.Name()

	meta, err = web.MockGet(app, fmt.Sprintf("/job.invocation/%s/%s", jobName, uuid.V4().String())).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusNotFound, meta.StatusCode)
}

func TestManagementServerPause(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()
	jm.StartAsync()
	defer jm.Stop()

	meta, err := web.MockGet(app, "/pause").Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.True(jm.Latch.IsPaused())
}

func TestManagementServerAPIPause(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()
	jm.StartAsync()
	defer jm.Stop()

	meta, err := web.MockPost(app, "/api/pause", nil).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.True(jm.Latch.IsPaused())
}

func TestManagementServerResume(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()
	jm.StartAsync()
	defer jm.Stop()

	jm.Pause()
	assert.True(jm.Latch.IsPaused())

	meta, err := web.MockGet(app, "/resume").Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)

	assert.True(jm.Latch.IsStarted())
}

func TestManagementServerAPIResume(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()
	jm.StartAsync()
	defer jm.Stop()

	jm.Pause()
	assert.True(jm.Latch.IsPaused())

	meta, err := web.MockPost(app, "/api/resume", nil).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.True(jm.Latch.IsStarted())
}

func TestManagementServerJobRun(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()

	job, err := jm.Job("test1")
	assert.Nil(err)
	assert.NotNil(job)
	jobName := job.Name()

	meta, err := web.MockGet(app, fmt.Sprintf("/job.run/%s", jobName)).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.NotNil(job.Last)
}

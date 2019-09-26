package jobkit

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/web"
)

func firstJob(jm *cron.JobManager) *cron.JobScheduler {
	sorted := sortedJobs(jm)
	if len(sorted) > 0 {
		return sorted[0]
	}
	return nil
}

// sortedJobs returns the list of jobs ordered by job name.
func sortedJobs(jm *cron.JobManager) []*cron.JobScheduler {
	var output []*cron.JobScheduler
	for _, js := range jm.Jobs {
		output = append(output, js)
	}
	sort.Sort(cron.JobSchedulersByJobNameAsc(output))
	return output
}

func createTestOutputChunk() cron.OutputChunk {
	return cron.OutputChunk{
		Timestamp: time.Now().UTC(),
		Data:      []byte(uuid.V4()),
	}
}

func createTestCompleteJobInvocation(jobName string, elapsed time.Duration) cron.JobInvocation {
	return cron.JobInvocation{
		ID:       uuid.V4().String(),
		JobName:  jobName,
		Started:  time.Now().UTC(),
		Finished: time.Now().UTC().Add(elapsed),
		State:    cron.JobInvocationStateComplete,
		Elapsed:  elapsed,
		Output: &cron.OutputBuffer{
			Chunks: []cron.OutputChunk{
				createTestOutputChunk(),
				createTestOutputChunk(),
				createTestOutputChunk(),
				createTestOutputChunk(),
				createTestOutputChunk(),
			},
		},
	}
}

func createTestFailedJobInvocation(jobName string, elapsed time.Duration, err error) cron.JobInvocation {
	return cron.JobInvocation{
		ID:       uuid.V4().String(),
		JobName:  jobName,
		Started:  time.Now().UTC(),
		Finished: time.Now().UTC().Add(elapsed),
		State:    cron.JobInvocationStateFailed,
		Elapsed:  elapsed,
		Err:      err,
		Output: &cron.OutputBuffer{
			Chunks: []cron.OutputChunk{
				createTestOutputChunk(),
				createTestOutputChunk(),
			},
		},
	}
}

func createTestJobManager() *cron.JobManager {
	test0 := cron.NewJob(cron.OptJobName("test0"))
	test1 := cron.NewJob(cron.OptJobName("test1"))
	test2 := cron.NewJob(cron.OptJobName("test2 job.foo"))

	jm := cron.New()
	jm.LoadJobs(test0, test1, test2)

	test0CurrentOutput := &cron.OutputBuffer{
		Chunks: []cron.OutputChunk{
			createTestOutputChunk(),
			createTestOutputChunk(),
			createTestOutputChunk(),
			createTestOutputChunk(),
		},
	}
	test0CurrentOutputListeners := new(cron.OutputListeners)
	test0CurrentOutput.Listener = test0CurrentOutputListeners.Trigger

	jm.Jobs["test0"].Current = &cron.JobInvocation{
		ID:              uuid.V4().String(),
		JobName:         "test0",
		Started:         time.Now().UTC(),
		Output:          test0CurrentOutput,
		OutputListeners: test0CurrentOutputListeners,
	}

	jm.Jobs["test0"].History = []cron.JobInvocation{
		createTestCompleteJobInvocation("test0", 200*time.Millisecond),
		createTestCompleteJobInvocation("test0", 250*time.Millisecond),
		createTestFailedJobInvocation("test0", 5*time.Second, fmt.Errorf("this is only a test %s", uuid.V4().String())),
	}
	jm.Jobs["test1"].History = []cron.JobInvocation{
		createTestCompleteJobInvocation("test1", 200*time.Millisecond),
		createTestCompleteJobInvocation("test1", 250*time.Millisecond),
		createTestCompleteJobInvocation("test1", 300*time.Millisecond),
		createTestCompleteJobInvocation("test1", 350*time.Millisecond),
	}
	jm.Jobs["test2 job.foo"].History = []cron.JobInvocation{
		createTestCompleteJobInvocation("test2 job.foo", 200*time.Millisecond),
		createTestCompleteJobInvocation("test2 job.foo", 250*time.Millisecond),
		createTestCompleteJobInvocation("test2 job.foo", 300*time.Millisecond),
		createTestCompleteJobInvocation("test2 job.foo", 350*time.Millisecond),
	}
	return jm
}

func createTestManagementServer() (*cron.JobManager, *web.App) {
	jm := createTestJobManager()
	return jm, NewServer(jm, Config{
		Web: web.Config{
			Port: 5000,
		},
	})
}

func TestMain(m *testing.M) {
	assert.Main(m)
}
